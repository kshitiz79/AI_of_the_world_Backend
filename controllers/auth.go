package controllers

import (
	"net/http"
	"time"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// Register creates a new user account
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if username exists
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Username already exists")
		return
	}

	// Check if email exists
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         "user",
		IsActive:     true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check if user is active
	if !user.IsActive {
		utils.ErrorResponse(c, http.StatusForbidden, "Account is deactivated")
		return
	}

	// Verify password
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	config.DB.Save(&user)

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Remove password hash from response
	user.PasswordHash = ""

	utils.SuccessResponse(c, http.StatusOK, "Login successful", models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	user.PasswordHash = ""
	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateInterests updates the user's interests
func UpdateInterests(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Interests string `json:"interests" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update user interests
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("interests", req.Interests).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update interests")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Interests updated successfully", nil)
}
