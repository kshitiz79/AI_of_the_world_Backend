package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"ai-of-the-world-backend/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() error {
	var err error
	cld, err = cloudinary.NewFromParams(
		config.AppConfig.CloudinaryCloudName,
		config.AppConfig.CloudinaryAPIKey,
		config.AppConfig.CloudinaryAPISecret,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}
	return nil
}

// UploadImage uploads an image to Cloudinary
func UploadImage(file multipart.File, filename string) (string, error) {
	if cld == nil {
		return "", fmt.Errorf("Cloudinary not initialized")
	}

	ctx := context.Background()

	// Generate unique filename
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	uniqueFilename := fmt.Sprintf("%s_%d%s", nameWithoutExt, time.Now().Unix(), ext)

	// Upload parameters
	uploadParams := uploader.UploadParams{
		Folder:         config.AppConfig.CloudinaryUploadFolder,
		PublicID:       strings.TrimSuffix(uniqueFilename, ext),
		ResourceType:   "image",
		Transformation: "q_auto,f_auto", // Auto quality and format
	}

	// Upload to Cloudinary
	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}

	return result.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary
func DeleteImage(publicID string) error {
	if cld == nil {
		return fmt.Errorf("Cloudinary not initialized")
	}

	ctx := context.Background()

	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete from Cloudinary: %w", err)
	}

	return nil
}

// GetPublicIDFromURL extracts the public ID from a Cloudinary URL
func GetPublicIDFromURL(url string) string {
	// Example URL: https://res.cloudinary.com/cloud_name/image/upload/v1234567890/folder/filename.jpg
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}

	// Get the last part (filename with extension)
	filename := parts[len(parts)-1]

	// Remove extension
	publicID := strings.TrimSuffix(filename, filepath.Ext(filename))

	// If there's a folder, include it
	if len(parts) >= 3 {
		folder := parts[len(parts)-2]
		if folder != "upload" && !strings.HasPrefix(folder, "v") {
			publicID = folder + "/" + publicID
		}
	}

	return publicID
}
