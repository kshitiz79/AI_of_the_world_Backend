package controllers

import (
	"net/http"
	"strings"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadGIF handles GIF upload to Backblaze B2
func UploadGIF(c *gin.Context) {
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

	// Get the GIF file
	file, fileHeader, err := c.Request.FormFile("gif")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "GIF file is required")
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/gif") {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only GIF files are allowed")
		return
	}

	// Upload to B2
	gifURL, err := utils.UploadToB2(file, fileHeader, "gifs", config.AppConfig.B2S3BucketGIF)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload GIF: "+err.Error())
		return
	}

	// Get form data
	projectTitle := c.PostForm("project_title")
	prompt := c.PostForm("prompt")
	technicalNotes := c.PostForm("technical_notes")
	modelOrTool := c.PostForm("model_or_tool")
	creatorCredit := c.PostForm("creator_credit")
	tagsStr := c.PostForm("tags")

	// Create GIF prompt record
	gifPrompt := models.GIFPrompt{
		UserID:         userID.(uint),
		ProjectTitle:   projectTitle,
		Prompt:         prompt,
		GIFURL:         gifURL,
		TechnicalNotes: technicalNotes,
		ModelOrTool:    modelOrTool,
		CreatorCredit:  creatorCredit,
		Status:         "pending",
		IsPublished:    false,
		IsFeatured:     false,
	}

	if err := config.DB.Create(&gifPrompt).Error; err != nil {
		// If database save fails, delete the uploaded file
		utils.DeleteFromB2(gifURL, config.AppConfig.B2S3BucketGIF)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save GIF prompt")
		return
	}

	// Handle tags if provided
	if tagsStr != "" {
		tagIDs := strings.Split(tagsStr, ",")
		for _, tagIDStr := range tagIDs {
			var tag models.Tag
			if err := config.DB.First(&tag, tagIDStr).Error; err == nil {
				config.DB.Model(&gifPrompt).Association("Tags").Append(&tag)
			}
		}
	}

	// Load relationships
	config.DB.Preload("User").Preload("Tags").First(&gifPrompt, gifPrompt.ID)

	utils.SuccessResponse(c, http.StatusCreated, "GIF uploaded successfully", gifPrompt)
}

// GetGIFPrompts returns all GIF prompts with optional filters
func GetGIFPrompts(c *gin.Context) {
	var gifs []models.GIFPrompt

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

	if err := query.Find(&gifs).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch GIF prompts")
		return
	}

	// Generate signed URLs for each GIF (valid for 7 days)
	for i := range gifs {
		if gifs[i].GIFURL != "" {
			signedURL, err := utils.GetSignedURL(gifs[i].GIFURL, config.AppConfig.B2S3BucketGIF)
			if err != nil {
				// Log error but continue
				println("Warning: Failed to generate signed URL for GIF", gifs[i].ID, ":", err.Error())
			} else {
				gifs[i].GIFURL = signedURL
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "GIF prompts retrieved successfully", gifs)
}

// GetGIFPromptByID returns a single GIF prompt by ID
func GetGIFPromptByID(c *gin.Context) {
	id := c.Param("id")

	var gif models.GIFPrompt
	if err := config.DB.Preload("User").Preload("Tags").First(&gif, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	// Generate signed URL (valid for 7 days)
	if gif.GIFURL != "" {
		signedURL, err := utils.GetSignedURL(gif.GIFURL, config.AppConfig.B2S3BucketGIF)
		if err != nil {
			println("Warning: Failed to generate signed URL for GIF", gif.ID, ":", err.Error())
		} else {
			gif.GIFURL = signedURL
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "GIF prompt retrieved successfully", gif)
}

// DeleteGIFPrompt deletes a GIF prompt (Admin or Owner)
func DeleteGIFPrompt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var prompt models.GIFPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	// Check if user is admin or owner
	if role != "admin" && prompt.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to delete this prompt")
		return
	}

	// Delete from B2
	if err := utils.DeleteFromB2(prompt.GIFURL, config.AppConfig.B2S3BucketGIF); err != nil {
		// Log error but continue with database deletion
		println("Warning: Failed to delete GIF from B2:", err.Error())
	}

	// Delete from database
	if err := config.DB.Delete(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete GIF prompt")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "GIF prompt deleted successfully", nil)
}

// ApproveGIFPrompt approves a pending GIF prompt (Admin only)
func ApproveGIFPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.GIFPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	prompt.Status = "approved"

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to approve GIF prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "GIF prompt approved successfully", prompt)
}

// RejectGIFPrompt rejects a pending GIF prompt (Admin only)
func RejectGIFPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.GIFPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	prompt.Status = "rejected"

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject GIF prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "GIF prompt rejected successfully", prompt)
}

// PublishGIFPrompt publishes an approved GIF prompt (Admin only)
func PublishGIFPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.GIFPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	if prompt.Status != "approved" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only approved prompts can be published")
		return
	}

	prompt.IsPublished = true

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to publish GIF prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "GIF prompt published successfully", prompt)
}

// UnpublishGIFPrompt unpublishes a GIF prompt (Admin only)
func UnpublishGIFPrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.GIFPrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "GIF prompt not found")
		return
	}

	prompt.IsPublished = false

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unpublish GIF prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "GIF prompt unpublished successfully", prompt)
}
