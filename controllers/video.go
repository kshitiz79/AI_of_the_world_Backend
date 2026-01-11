package controllers

import (
	"net/http"
	"strings"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadVideo handles video upload to Backblaze B2
func UploadVideo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(config.AppConfig.MaxUploadSize); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File too large or invalid")
		return
	}

	// Get the video file
	file, fileHeader, err := c.Request.FormFile("video")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Video file is required")
		return
	}
	defer file.Close()

	// Validate file type (video)
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "video/") {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only video files are allowed")
		return
	}

	// Upload to B2
	videoURL, err := utils.UploadToB2(file, fileHeader, "videos", config.AppConfig.B2S3BucketVideo)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload video: "+err.Error())
		return
	}

	// Get form data
	projectTitle := c.PostForm("project_title")
	prompt := c.PostForm("prompt")
	technicalNotes := c.PostForm("technical_notes")
	modelOrTool := c.PostForm("model_or_tool")
	creatorCredit := c.PostForm("creator_credit")
	tagsStr := c.PostForm("tags")

	// Create video prompt record
	videoPrompt := models.VideoPrompt{
		UserID:         userID.(uint),
		ProjectTitle:   projectTitle,
		Prompt:         prompt,
		VideoURL:       videoURL,
		TechnicalNotes: technicalNotes,
		ModelOrTool:    modelOrTool,
		CreatorCredit:  creatorCredit,
		Status:         "pending",
		IsPublished:    false,
		IsFeatured:     false,
	}

	if err := config.DB.Create(&videoPrompt).Error; err != nil {
		// If database save fails, delete the uploaded file
		utils.DeleteFromB2(videoURL, config.AppConfig.B2S3BucketVideo)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save video prompt")
		return
	}

	// Handle tags if provided
	if tagsStr != "" {
		tagIDs := strings.Split(tagsStr, ",")
		for _, tagIDStr := range tagIDs {
			var tag models.Tag
			if err := config.DB.First(&tag, tagIDStr).Error; err == nil {
				config.DB.Model(&videoPrompt).Association("Tags").Append(&tag)
			}
		}
	}

	// Load relationships
	config.DB.Preload("User").Preload("Tags").First(&videoPrompt, videoPrompt.ID)

	utils.SuccessResponse(c, http.StatusCreated, "Video uploaded successfully", videoPrompt)
}

// GetVideoPrompts returns all video prompts with optional filters
func GetVideoPrompts(c *gin.Context) {
	var videos []models.VideoPrompt

	query := config.DB.Preload("User").Preload("Tags")

	// Optional filters
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	userID := c.Query("user_id")
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	isFeatured := c.Query("is_featured")
	if isFeatured == "true" {
		query = query.Where("is_featured = ?", true)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&videos).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch video prompts")
		return
	}

	// Generate signed URLs for each video (valid for 7 days)
	for i := range videos {
		if videos[i].VideoURL != "" {
			signedURL, err := utils.GetSignedURL(videos[i].VideoURL, config.AppConfig.B2S3BucketVideo)
			if err != nil {
				// Log error but continue
				println("Warning: Failed to generate signed URL for video", videos[i].ID, ":", err.Error())
			} else {
				videos[i].VideoURL = signedURL
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Video prompts retrieved successfully", videos)
}

// GetVideoPromptByID returns a single video prompt by ID
func GetVideoPromptByID(c *gin.Context) {
	id := c.Param("id")

	var video models.VideoPrompt
	if err := config.DB.Preload("User").Preload("Tags").First(&video, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	// Generate signed URL (valid for 7 days)
	if video.VideoURL != "" {
		signedURL, err := utils.GetSignedURL(video.VideoURL, config.AppConfig.B2S3BucketVideo)
		if err != nil {
			println("Warning: Failed to generate signed URL for video", video.ID, ":", err.Error())
		} else {
			video.VideoURL = signedURL
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Video prompt retrieved successfully", video)
}

// DeleteVideoPrompt deletes a video prompt (Admin or Owner)
func DeleteVideoPrompt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var prompt models.VideoPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	// Check if user is admin or owner
	if role != "admin" && prompt.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to delete this prompt")
		return
	}

	// Delete from B2
	if err := utils.DeleteFromB2(prompt.VideoURL, config.AppConfig.B2S3BucketVideo); err != nil {
		// Log error but continue with database deletion
		println("Warning: Failed to delete video from B2:", err.Error())
	}

	// Delete from database
	if err := config.DB.Delete(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete video prompt")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Video prompt deleted successfully", nil)
}

// ApproveVideoPrompt approves a pending video prompt (Admin only)
func ApproveVideoPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.VideoPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	prompt.Status = "approved"

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to approve video prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Video prompt approved successfully", prompt)
}

// RejectVideoPrompt rejects a pending video prompt (Admin only)
func RejectVideoPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.VideoPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	prompt.Status = "rejected"

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject video prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Video prompt rejected successfully", prompt)
}

// PublishVideoPrompt publishes an approved video prompt (Admin only)
func PublishVideoPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.VideoPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	if prompt.Status != "approved" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only approved prompts can be published")
		return
	}

	prompt.IsPublished = true

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to publish video prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Video prompt published successfully", prompt)
}

// UnpublishVideoPrompt unpublishes a video prompt (Admin only)
func UnpublishVideoPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.VideoPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Video prompt not found")
		return
	}

	prompt.IsPublished = false

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unpublish video prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Video prompt unpublished successfully", prompt)
}
