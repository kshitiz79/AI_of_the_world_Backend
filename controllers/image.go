package controllers

import (
	"net/http"
	"strings"
	"time"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// UploadImage handles image upload to Cloudinary
func UploadImage(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse multipart form
	err := c.Request.ParseMultipartForm(config.AppConfig.MaxUploadSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File too large or invalid")
		return
	}

	// Get the file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No image file provided")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		utils.ErrorResponse(c, http.StatusBadRequest, "File must be an image")
		return
	}

	// Upload to Cloudinary
	imageURL, err := utils.UploadImage(file, header.Filename)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image: "+err.Error())
		return
	}

	// Get form data
	projectTitle := c.PostForm("project_title")
	prompt := c.PostForm("prompt")
	technicalNotes := c.PostForm("technical_notes")
	modelOrTool := c.PostForm("model_or_tool")
	creatorCredit := c.PostForm("creator_credit")
	tagsStr := c.PostForm("tags") // Comma-separated tag IDs

	// Validate required fields
	if projectTitle == "" || prompt == "" || creatorCredit == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Create image prompt
	imagePrompt := models.ImagePrompt{
		UserID:         userID.(uint),
		ProjectTitle:   projectTitle,
		Prompt:         prompt,
		TechnicalNotes: technicalNotes,
		ModelOrTool:    modelOrTool,
		CreatorCredit:  creatorCredit,
		ImageURL:       imageURL,
		ImageFilename:  header.Filename,
		Status:         "pending",
		IsPublished:    false,
	}

	// Save to database
	if err := config.DB.Create(&imagePrompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save image prompt")
		return
	}

	// Handle tags if provided
	if tagsStr != "" {
		tagIDs := strings.Split(tagsStr, ",")
		var tags []models.Tag

		if err := config.DB.Where("id IN ?", tagIDs).Find(&tags).Error; err == nil {
			config.DB.Model(&imagePrompt).Association("Tags").Append(tags)
		}
	}

	// Load the prompt with user and tags
	config.DB.Preload("User").Preload("Tags").First(&imagePrompt, imagePrompt.ID)

	utils.SuccessResponse(c, http.StatusCreated, "Image uploaded successfully", imagePrompt)
}

// GetImagePrompts returns all image prompts (with filters)
func GetImagePrompts(c *gin.Context) {
	var prompts []models.ImagePrompt

	query := config.DB.Preload("User").Preload("Tags")

	// Filters
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

	if err := query.Find(&prompts).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch image prompts")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image prompts retrieved successfully", prompts)
}

// GetImagePromptByID returns a single image prompt
func GetImagePromptByID(c *gin.Context) {
	id := c.Param("id")

	var prompt models.ImagePrompt
	if err := config.DB.Preload("User").Preload("Tags").First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image prompt retrieved successfully", prompt)
}

// DeleteImagePrompt deletes an image prompt (Admin or Owner)
func DeleteImagePrompt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	// Check if user is admin or owner
	if role != "admin" && prompt.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to delete this prompt")
		return
	}

	// Delete from Cloudinary
	publicID := utils.GetPublicIDFromURL(prompt.ImageURL)
	if publicID != "" {
		utils.DeleteImage(publicID)
	}

	// Delete from database
	if err := config.DB.Delete(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete image prompt")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Image prompt deleted successfully", nil)
}

// ApproveImagePrompt approves a pending image prompt (Admin only)
func ApproveImagePrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	prompt.Status = "approved"
	now := time.Now()
	prompt.VerifiedAt = &now

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to approve image prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Image prompt approved successfully", prompt)
}

// RejectImagePrompt rejects a pending image prompt (Admin only)
func RejectImagePrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	prompt.Status = "rejected"

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject image prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Image prompt rejected successfully", prompt)
}

// PublishImagePrompt publishes an approved image prompt (Admin only)
func PublishImagePrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	if prompt.Status != "approved" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Only approved prompts can be published")
		return
	}

	prompt.IsPublished = true

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to publish image prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Image prompt published successfully", prompt)
}

// UnpublishImagePrompt unpublishes an image prompt (Admin only)
func UnpublishImagePrompt(c *gin.Context) {
	id := c.Param("id")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	prompt.IsPublished = false

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unpublish image prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Image prompt unpublished successfully", prompt)
}

// UpdateImagePrompt updates an image prompt (Admin or Owner)
func UpdateImagePrompt(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var prompt models.ImagePrompt
	if err := config.DB.First(&prompt, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Image prompt not found")
		return
	}

	// Check if user is admin or owner
	if role != "admin" && prompt.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "You don't have permission to update this prompt")
		return
	}

	var req struct {
		ProjectTitle   string `json:"project_title"`
		Prompt         string `json:"prompt"`
		TechnicalNotes string `json:"technical_notes"`
		ModelOrTool    string `json:"model_or_tool"`
		CreatorCredit  string `json:"creator_credit"`
		IsFeatured     *bool  `json:"is_featured"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.ProjectTitle != "" {
		prompt.ProjectTitle = req.ProjectTitle
	}
	if req.Prompt != "" {
		prompt.Prompt = req.Prompt
	}
	if req.TechnicalNotes != "" {
		prompt.TechnicalNotes = req.TechnicalNotes
	}
	if req.ModelOrTool != "" {
		prompt.ModelOrTool = req.ModelOrTool
	}
	if req.CreatorCredit != "" {
		prompt.CreatorCredit = req.CreatorCredit
	}
	if req.IsFeatured != nil && role == "admin" {
		prompt.IsFeatured = *req.IsFeatured
	}

	if err := config.DB.Save(&prompt).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update image prompt")
		return
	}

	config.DB.Preload("User").Preload("Tags").First(&prompt, prompt.ID)
	utils.SuccessResponse(c, http.StatusOK, "Image prompt updated successfully", prompt)
}
