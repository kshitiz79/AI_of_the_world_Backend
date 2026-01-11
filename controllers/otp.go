package controllers

import (
	"net/http"
	"time"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/models"
	"ai-of-the-world-backend/utils"

	"github.com/gin-gonic/gin"
)

// SendOTP sends an OTP to the user's email
func SendOTP(c *gin.Context) {
	var req models.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// For signup, check if email already exists
	if req.Purpose == "signup" {
		var existingUser models.User
		if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
			return
		}
	}

	// For forgot_password, check if email exists
	if req.Purpose == "forgot_password" {
		var user models.User
		if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "Email not found")
			return
		}
	}

	// Generate OTP
	otp, err := utils.GenerateOTP()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate OTP")
		return
	}

	// Delete any existing unverified OTPs for this email and purpose
	config.DB.Where("email = ? AND purpose = ? AND verified = ?", req.Email, req.Purpose, false).Delete(&models.OTP{})

	// Save OTP to database
	otpRecord := models.OTP{
		Email:     req.Email,
		OTP:       otp,
		Purpose:   req.Purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Verified:  false,
	}

	if err := config.DB.Create(&otpRecord).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save OTP")
		return
	}

	// Send OTP via email
	if err := utils.SendOTPEmail(req.Email, otp, req.Purpose); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to send OTP email: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP sent successfully", gin.H{
		"email":      req.Email,
		"expires_at": otpRecord.ExpiresAt,
	})
}

// VerifyOTP verifies the OTP without completing signup/reset
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find OTP
	var otpRecord models.OTP
	if err := config.DB.Where("email = ? AND otp = ? AND verified = ?", req.Email, req.OTP, false).
		Order("created_at DESC").First(&otpRecord).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid OTP")
		return
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "OTP has expired")
		return
	}

	// Mark OTP as verified
	otpRecord.Verified = true
	config.DB.Save(&otpRecord)

	utils.SuccessResponse(c, http.StatusOK, "OTP verified successfully", nil)
}

// SignupWithOTP completes the signup process with OTP verification
func SignupWithOTP(c *gin.Context) {
	var req models.SignupWithOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify OTP
	var otpRecord models.OTP
	if err := config.DB.Where("email = ? AND otp = ? AND purpose = ? AND verified = ?",
		req.Email, req.OTP, "signup", true).
		Order("created_at DESC").First(&otpRecord).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or unverified OTP")
		return
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "OTP has expired")
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
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered")
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
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		FullName:      req.FullName,
		Role:          "user",
		IsActive:      true,
		EmailVerified: true, // Email is verified via OTP
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Delete used OTP
	config.DB.Delete(&otpRecord)

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

// ResetPassword resets the user's password using OTP
func ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify OTP
	var otpRecord models.OTP
	if err := config.DB.Where("email = ? AND otp = ? AND purpose = ? AND verified = ?",
		req.Email, req.OTP, "forgot_password", true).
		Order("created_at DESC").First(&otpRecord).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or unverified OTP")
		return
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "OTP has expired")
		return
	}

	// Find user
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	user.PasswordHash = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Delete used OTP
	config.DB.Delete(&otpRecord)

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}
