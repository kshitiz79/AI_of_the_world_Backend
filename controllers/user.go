package controllers

import (
	"net/http"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// GetAllUsers returns all users (Admin only)
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Optional filters
	role := c.Query("role")
	isActive := c.Query("is_active")

	query := config.DB

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&users).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	// Remove password hashes from response
	for i := range users {
		users[i].PasswordHash = ""
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", users)
}

// GetUserByID returns a single user by ID (Admin only)
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUserStatus updates user's active status (Admin only)
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.IsActive = req.IsActive

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "User status updated successfully", user)
}

// DeleteUser deletes a user (Admin only)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Prevent deleting yourself
	currentUserID, _ := c.Get("userID")
	if user.ID == currentUserID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "Cannot delete your own account")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// GetUserStats returns user statistics (Admin only)
func GetUserStats(c *gin.Context) {
	var totalUsers int64
	var totalAdmins int64
	var activeUsers int64

	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&totalAdmins)
	config.DB.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)

	stats := map[string]interface{}{
		"total_users":   totalUsers,
		"total_admins":  totalAdmins,
		"active_users":  activeUsers,
		"regular_users": totalUsers - totalAdmins,
	}

	utils.SuccessResponse(c, http.StatusOK, "User statistics retrieved successfully", stats)
}
