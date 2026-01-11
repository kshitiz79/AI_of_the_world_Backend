package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SendOTPEmail sends an OTP to the specified email address
func SendOTPEmail(toEmail, otp, purpose string) error {
	// Get email configuration from environment variables
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	fromEmail := os.Getenv("SMTP_EMAIL")
	appPassword := os.Getenv("SMTP_PASSWORD")

	if fromEmail == "" || appPassword == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	// Determine email subject and body based on purpose
	var subject, body string
	switch purpose {
	case "signup":
		subject = "Verify Your Email - AI of the World"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #e5e7eb; background-color: #000000; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #000000 0%%, #1f2937 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; border: 1px solid #374151; }
        .content { background: #111827; padding: 30px; border-radius: 0 0 10px 10px; border: 1px solid #374151; border-top: none; }
        .otp-box { background: #000000; border: 2px dashed #06b6d4; padding: 16px; text-align: center; margin: 20px 0; border-radius: 8px; }
        .otp-code { font-size: 22px; font-weight: bold; color: #06b6d4; letter-spacing: 3px; }
        .footer { text-align: center; margin-top: 20px; color: #9ca3af; font-size: 12px; }
        h1 { margin: 0; font-size: 28px; }
        h2 { color: #f9fafb; margin-top: 0; }
        p { color: #d1d5db; }
        strong { color: #fbbf24; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h3> Welcome to AI of the World!</h1>
        </div>
        <div class="content">
            <h2>Email Verification</h2>
            <p>Thank you for signing up! Please use the following OTP to verify your email address:</p>
            <div class="otp-box">
                <div class="otp-code">%s</div>
            </div>
            <p><strong>This OTP will expire in 10 minutes.</strong></p>
            <p>If you didn't request this verification, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>© 2026 AI of the World. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, otp)
	case "forgot_password":
		subject = "Reset Your Password - AI of the World"
		body = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #e5e7eb; background-color: #000000; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #000000 0%%, #1f2937 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; border: 1px solid #374151; }
        .content { background: #111827; padding: 30px; border-radius: 0 0 10px 10px; border: 1px solid #374151; border-top: none; }
        .otp-box { background: #000000; border: 2px dashed #06b6d4; padding: 20px; text-align: center; margin: 20px 0; border-radius: 8px; }
        .otp-code { font-size: 32px; font-weight: bold; color: #06b6d4; letter-spacing: 8px; }
        .footer { text-align: center; margin-top: 20px; color: #9ca3af; font-size: 12px; }
        h1 { margin: 0; font-size: 28px; }
        h2 { color: #f9fafb; margin-top: 0; }
        p { color: #d1d5db; }
        strong { color: #fbbf24; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Reset Your Password</h2>
            <p>We received a request to reset your password. Please use the following OTP to proceed:</p>
            <div class="otp-box">
                <div class="otp-code">%s</div>
            </div>
            <p><strong>This OTP will expire in 10 minutes.</strong></p>
            <p>If you didn't request a password reset, please ignore this email and your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>© 2026 AI of the World. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, otp)
	default:
		return fmt.Errorf("invalid email purpose")
	}

	// Compose email message
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s\r\n", fromEmail, toEmail, subject, body))

	// Authentication
	auth := smtp.PlainAuth("", fromEmail, appPassword, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
