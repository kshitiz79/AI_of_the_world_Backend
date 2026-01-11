package models

import (
	"time"
)

// OTP represents an OTP record in the database
type OTP struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"size:255;not null;index" json:"email"`
	OTP       string    `gorm:"size:6;not null" json:"-"`
	Purpose   string    `gorm:"type:enum('signup','forgot_password');not null" json:"purpose"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	CreatedAt time.Time `json:"created_at"`
}

func (OTP) TableName() string {
	return "otps"
}

// SendOTPRequest represents the request to send an OTP
type SendOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required,oneof=signup forgot_password"`
}

// VerifyOTPRequest represents the request to verify an OTP
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// SignupWithOTPRequest represents the complete signup request with OTP
type SignupWithOTPRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	OTP      string `json:"otp" binding:"required,len=6"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
