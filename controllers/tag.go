package controllers

import (
	"net/http"
	"strconv"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// GetAllTags returns all tags
func GetAllTags(c *gin.Context) {
	var tags []models.Tag

	// Optional filters
	category := c.Query("category")
	isActive := c.Query("is_active")

	query := config.DB

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if err := query.Order("usage_count DESC, name ASC").Find(&tags).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch tags")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tags retrieved successfully", tags)
}

// GetTagByID returns a single tag by ID
func GetTagByID(c *gin.Context) {
	id := c.Param("id")

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Tag not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tag retrieved successfully", tag)
}

// CreateTag creates a new tag (Admin only)
func CreateTag(c *gin.Context) {
	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if tag name already exists
	var existingTag models.Tag
	if err := config.DB.Where("name = ?", req.Name).First(&existingTag).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Tag with this name already exists")
		return
	}

	tag := models.Tag{
		Name:        req.Name,
		Category:    req.Category,
		Description: req.Description,
		IsActive:    true,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create tag")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Tag created successfully", tag)
}

// UpdateTag updates an existing tag (Admin only)
func UpdateTag(c *gin.Context) {
	id := c.Param("id")

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Tag not found")
		return
	}

	var req models.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if new name conflicts with existing tag
	if req.Name != "" && req.Name != tag.Name {
		var existingTag models.Tag
		if err := config.DB.Where("name = ? AND id != ?", req.Name, id).First(&existingTag).Error; err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Tag with this name already exists")
			return
		}
		tag.Name = req.Name
	}

	if req.Category != "" {
		tag.Category = req.Category
	}

	if req.Description != "" {
		tag.Description = req.Description
	}

	if req.IsActive != nil {
		tag.IsActive = *req.IsActive
	}

	if err := config.DB.Save(&tag).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update tag")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tag updated successfully", tag)
}

// DeleteTag deletes a tag (Admin only)
func DeleteTag(c *gin.Context) {
	id := c.Param("id")

	var tag models.Tag
	if err := config.DB.First(&tag, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Tag not found")
		return
	}

	if err := config.DB.Delete(&tag).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete tag")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tag deleted successfully", nil)
}

// GetTagStats returns statistics about tags
func GetTagStats(c *gin.Context) {
	var totalTags int64
	config.DB.Model(&models.Tag{}).Count(&totalTags)

	var categoryStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	config.DB.Model(&models.Tag{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Scan(&categoryStats)

	stats := map[string]interface{}{
		"total_tags":     totalTags,
		"category_stats": categoryStats,
	}

	utils.SuccessResponse(c, http.StatusOK, "Tag statistics retrieved successfully", stats)
}

// SearchTags searches tags by name
func SearchTags(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Search query required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var tags []models.Tag
	if err := config.DB.Where("name LIKE ?", "%"+query+"%").
		Limit(limit).
		Find(&tags).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search tags")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Search results", tags)
}
