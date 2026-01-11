# Email OTP Setup Guide

## Overview
This application now supports email-based OTP (One-Time Password) verification for:
- User signup with email verification
- Password reset functionality

## SMTP Configuration

### Gmail App Password Setup
The application uses Gmail's SMTP server to send OTP emails. You need to add the following environment variables to your `.env` file:

```env
# Email Configuration (Gmail SMTP)
SMTP_EMAIL=your-gmail@gmail.com
SMTP_PASSWORD=yefb iude pmjo askn
```

**Important Notes:**
1. The `SMTP_PASSWORD` is a **Google App Password**, not your regular Gmail password
2. The provided app password is: `yefb iude pmjo askn`
3. Make sure to use the email address associated with this app password

### How to Generate a New Google App Password (if needed)
1. Go to your Google Account settings
2. Navigate to Security → 2-Step Verification
3. Scroll down to "App passwords"
4. Generate a new app password for "Mail"
5. Copy the 16-character password (spaces will be added automatically)
6. Add it to your `.env` file

## API Endpoints

### 1. Send OTP
**Endpoint:** `POST /api/v1/auth/send-otp`

**Request Body:**
```json
{
  "email": "user@example.com",
  "purpose": "signup" // or "forgot_password"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "email": "user@example.com",
    "expires_at": "2026-01-12T00:27:51Z"
  }
}
```

### 2. Verify OTP
**Endpoint:** `POST /api/v1/auth/verify-otp`

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

**Response:**
```json
{
  "success": true,
  "message": "OTP verified successfully"
}
```

### 3. Signup with OTP
**Endpoint:** `POST /api/v1/auth/signup-with-otp`

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456",
  "username": "johndoe",
  "full_name": "John Doe",
  "password": "SecurePass123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "user@example.com",
      "full_name": "John Doe",
      "role": "user",
      "email_verified": true,
      "is_active": true
    }
  }
}
```

### 4. Reset Password
**Endpoint:** `POST /api/v1/auth/reset-password`

**Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456",
  "new_password": "NewSecurePass123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

## Frontend Pages

### 1. Signup Page
**URL:** `http://localhost:3000/signup`

Features:
- Multi-step signup process
- Email verification via OTP
- Password strength validation
- Real-time error handling

### 2. Forgot Password Page
**URL:** `http://localhost:3000/forgot-password`

Features:
- Email-based password reset
- OTP verification
- New password creation
- Automatic redirect to login after success

### 3. Dashboard Protection
The user dashboard at `http://localhost:3000/panel/your-dashboard/` is now protected with `AuthGuard`. Users must be logged in to access it. Unauthenticated users will be redirected to the signin page.

## OTP Details

- **OTP Length:** 6 digits
- **Expiration Time:** 10 minutes
- **Email Template:** HTML-formatted with gradient design
- **Resend Capability:** Users can request a new OTP if needed

## Database Schema

A new `otps` table has been created with the following structure:

```sql
CREATE TABLE otps (
  id INT PRIMARY KEY AUTO_INCREMENT,
  email VARCHAR(255) NOT NULL,
  otp VARCHAR(6) NOT NULL,
  purpose ENUM('signup', 'forgot_password') NOT NULL,
  expires_at DATETIME NOT NULL,
  verified BOOLEAN DEFAULT FALSE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_email (email)
);
```

## Testing

### Test Signup Flow
1. Navigate to `http://localhost:3000/signup`
2. Enter your email address
3. Check your email for the OTP
4. Enter the 6-digit OTP
5. Complete your profile with username, full name, and password
6. You'll be automatically logged in and redirected to your dashboard

### Test Forgot Password Flow
1. Navigate to `http://localhost:3000/signin`
2. Click "Forgot password?"
3. Enter your registered email
4. Check your email for the OTP
5. Enter the OTP and your new password
6. You'll be redirected to the signin page

## Security Features

1. **OTP Expiration:** OTPs expire after 10 minutes
2. **One-Time Use:** OTPs are deleted after successful verification
3. **Email Verification:** Users' emails are marked as verified after OTP confirmation
4. **Password Hashing:** All passwords are securely hashed using bcrypt
5. **JWT Authentication:** Secure token-based authentication
6. **Protected Routes:** Dashboard routes require valid authentication

## Troubleshooting

### Email Not Sending
- Verify SMTP credentials in `.env` file
- Check if 2-Step Verification is enabled on your Google account
- Ensure the app password is correct (no extra spaces)
- Check backend logs for SMTP errors

### OTP Not Working
- Verify OTP hasn't expired (10-minute limit)
- Ensure you're using the latest OTP sent to your email
- Check that the email address matches exactly

### Dashboard Access Issues
- Clear browser localStorage and try logging in again
- Verify the JWT token is being stored correctly
- Check browser console for authentication errors

## Next Steps

To enable this feature in production:
1. Add the SMTP credentials to your production `.env` file
2. Consider using a dedicated email service (SendGrid, AWS SES, etc.) for better deliverability
3. Implement rate limiting for OTP requests to prevent abuse
4. Add email templates for different languages if needed
5. Monitor email delivery rates and OTP usage
