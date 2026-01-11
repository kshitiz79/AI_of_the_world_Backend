package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"ai-of-the-world-backend/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var B2Session *session.Session
var B2Uploader *s3manager.Uploader
var B2Service *s3.S3

// InitializeB2 initializes the Backblaze B2 S3-compatible client
func InitializeB2() error {
	cfg := config.AppConfig

	if cfg.B2S3AccessKey == "" || cfg.B2S3SecretKey == "" {
		return fmt.Errorf("B2 S3 credentials not configured")
	}

	// Create AWS session with B2 credentials
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.B2S3Region),
		Endpoint:         aws.String(fmt.Sprintf("https://%s", cfg.B2S3Endpoint)),
		Credentials:      credentials.NewStaticCredentials(cfg.B2S3AccessKey, cfg.B2S3SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return fmt.Errorf("failed to create B2 session: %v", err)
	}

	B2Session = sess
	B2Uploader = s3manager.NewUploader(sess)
	B2Service = s3.New(sess)

	return nil
}

// UploadToB2 uploads a file to Backblaze B2
func UploadToB2(file multipart.File, fileHeader *multipart.FileHeader, folder string, bucketName string) (string, error) {
	if B2Uploader == nil {
		return "", fmt.Errorf("B2 uploader not initialized")
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s/%d%s", folder, time.Now().UnixNano(), ext)

	// Read file content
	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	file.Seek(0, 0) // Reset file pointer

	// Determine content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to B2
	result, err := B2Uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to B2: %v", err)
	}

	return result.Location, nil
}

// DeleteFromB2 deletes a file from Backblaze B2
func DeleteFromB2(fileURL string, bucketName string) error {
	if B2Service == nil {
		return fmt.Errorf("B2 service not initialized")
	}

	// Extract key from URL
	key := GetB2KeyFromURL(fileURL)
	if key == "" {
		return fmt.Errorf("invalid B2 URL")
	}

	_, err := B2Service.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from B2: %v", err)
	}

	return nil
}

// GetB2KeyFromURL extracts the object key from a B2 URL
func GetB2KeyFromURL(url string) string {
	// B2 URL format: https://s3.us-east-005.backblazeb2.com/bucket-name/folder/file.ext
	parts := strings.Split(url, "/")
	if len(parts) < 5 {
		return ""
	}
	// Join everything after the bucket name
	return strings.Join(parts[4:], "/")
}

// GetSignedURL generates a pre-signed URL for private B2 files
// URLs are valid for 7 days (maximum allowed by AWS S3/B2)
func GetSignedURL(fileURL string, bucketName string) (string, error) {
	if B2Service == nil {
		return "", fmt.Errorf("B2 service not initialized")
	}

	// Extract key from URL
	key := GetB2KeyFromURL(fileURL)
	if key == "" {
		return "", fmt.Errorf("invalid B2 URL")
	}

	// Generate pre-signed URL
	req, _ := B2Service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	// Set expiration time to 7 days (maximum allowed by AWS S3/B2)
	// 7 days = 604800 seconds
	urlStr, err := req.Presign(7 * 24 * time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return urlStr, nil
}
