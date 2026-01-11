package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"ai-of-the-world-backend/config"
	"ai-of-the-world-backend/routes"
	"ai-of-the-world-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()
	defer config.CloseDatabase()

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		log.Println("⚠️  Cloudinary initialization failed:", err)
		log.Println("   Image uploads will not work. Please check your Cloudinary credentials.")
	} else {
		log.Println("✅ Cloudinary initialized successfully")
	}

	// Initialize Backblaze B2
	if err := utils.InitializeB2(); err != nil {
		log.Println("⚠️  Backblaze B2 initialization failed:", err)
		log.Println("   GIF uploads will not work. Please check your B2 credentials.")
	} else {
		log.Println("✅ Backblaze B2 initialized successfully")
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(config.AppConfig.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create uploads directory:", err)
	}

	// Set Gin mode
	if config.AppConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins:     config.AppConfig.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "AI of the World API is running",
		})
	})

	// Setup routes
	routes.SetupRoutes(router)

	// Serve uploaded files
	router.Static("/uploads", config.AppConfig.UploadDir)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("\n🛑 Shutting down server...")
		config.CloseDatabase()
		os.Exit(0)
	}()

	// Start server
	serverAddr := ":" + config.AppConfig.Port
	log.Printf("🚀 Server starting on http://localhost%s\n", serverAddr)
	log.Printf("📚 API Documentation: http://localhost%s/api/v1\n", serverAddr)
	log.Printf("💚 Health Check: http://localhost%s/health\n", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
