package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	Environment    string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	UploadDir      string
	MaxUploadSize  int64
	AllowedOrigins []string
	FrontendURL    string
	// Cloudinary
	CloudinaryCloudName    string
	CloudinaryAPIKey       string
	CloudinaryAPISecret    string
	CloudinaryUploadFolder string
	// Backblaze B2 S3
	B2S3AccessKey   string
	B2S3SecretKey   string
	B2S3Endpoint    string
	B2S3Region      string
	B2S3BucketGIF   string
	B2S3BucketVideo string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "104857600"), 10, 64)

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENV", "development"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "ai_of_the_world"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-this"),
		UploadDir:      getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize:  maxUploadSize,
		AllowedOrigins: strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		FrontendURL:    getEnv("FRONTEND_URL", "http://localhost:3000"),
		// Cloudinary
		CloudinaryCloudName:    getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:       getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:    getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryUploadFolder: getEnv("CLOUDINARY_UPLOAD_FOLDER", "uploads"),
		// Backblaze B2 S3
		B2S3AccessKey:   getEnv("B2_S3_ACCESS_KEY", ""),
		B2S3SecretKey:   getEnv("B2_S3_SECRET_KEY", ""),
		B2S3Endpoint:    getEnv("B2_S3_ENDPOINT", ""),
		B2S3Region:      getEnv("B2_S3_REGION", "us-east-005"),
		B2S3BucketGIF:   getEnv("B2_S3_BUCKET_GIF", "aiofhtheworlsgif"),
		B2S3BucketVideo: getEnv("B2_S3_BUCKET_VIDEO", "aiofhtheworlsvideo"),
	}

	log.Println("✅ Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
