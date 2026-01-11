package routes

import (
	"ai-of-the-world-backend/controllers"
	"ai-of-the-world-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)

			// OTP routes
			auth.POST("/send-otp", controllers.SendOTP)
			auth.POST("/verify-otp", controllers.VerifyOTP)
			auth.POST("/signup-with-otp", controllers.SignupWithOTP)
			auth.POST("/reset-password", controllers.ResetPassword)
		}

		// Public tag routes (read-only)
		tags := v1.Group("/tags")
		{
			tags.GET("", controllers.GetAllTags)
			tags.GET("/:id", controllers.GetTagByID)
			tags.GET("/search", controllers.SearchTags)
			tags.GET("/stats", controllers.GetTagStats)
		}

		// Public image prompts (read-only)
		images := v1.Group("/images")
		{
			images.GET("", controllers.GetImagePrompts)
			images.GET("/:id", controllers.GetImagePromptByID)
		}

		// Public GIF prompts (read-only)
		gifs := v1.Group("/gifs")
		{
			gifs.GET("", controllers.GetGIFPrompts)
			gifs.GET("/:id", controllers.GetGIFPromptByID)
		}

		// Public video prompts (read-only)
		videos := v1.Group("/videos")
		{
			videos.GET("", controllers.GetVideoPrompts)
			videos.GET("/:id", controllers.GetVideoPromptByID)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", controllers.GetProfile)
			protected.PUT("/profile/interests", controllers.UpdateInterests)

			// Image upload
			protected.POST("/images/upload", controllers.UploadImage)
			protected.DELETE("/images/:id", controllers.DeleteImagePrompt)

			// GIF upload
			protected.POST("/gifs/upload", controllers.UploadGIF)
			protected.DELETE("/gifs/:id", controllers.DeleteGIFPrompt)

			// Video upload
			protected.POST("/videos/upload", controllers.UploadVideo)
			protected.DELETE("/videos/:id", controllers.DeleteVideoPrompt)

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// Tag management
				admin.POST("/tags", controllers.CreateTag)
				admin.PUT("/tags/:id", controllers.UpdateTag)
				admin.DELETE("/tags/:id", controllers.DeleteTag)

				// User management
				admin.GET("/users/stats", controllers.GetUserStats)
				admin.GET("/users", controllers.GetAllUsers)
				admin.GET("/users/:id", controllers.GetUserByID)
				admin.PUT("/users/:id/status", controllers.UpdateUserStatus)
				admin.DELETE("/users/:id", controllers.DeleteUser)

				// Image management
				admin.PUT("/images/:id/approve", controllers.ApproveImagePrompt)
				admin.PUT("/images/:id/reject", controllers.RejectImagePrompt)
				admin.PUT("/images/:id/publish", controllers.PublishImagePrompt)
				admin.PUT("/images/:id/unpublish", controllers.UnpublishImagePrompt)
				admin.PUT("/images/:id", controllers.UpdateImagePrompt)

				// GIF management
				admin.PUT("/gifs/:id/approve", controllers.ApproveGIFPrompt)
				admin.PUT("/gifs/:id/reject", controllers.RejectGIFPrompt)
				admin.PUT("/gifs/:id/publish", controllers.PublishGIFPrompt)
				admin.PUT("/gifs/:id/unpublish", controllers.UnpublishGIFPrompt)

				// Video management
				admin.PUT("/videos/:id/approve", controllers.ApproveVideoPrompt)
				admin.PUT("/videos/:id/reject", controllers.RejectVideoPrompt)
				admin.PUT("/videos/:id/publish", controllers.PublishVideoPrompt)
				admin.PUT("/videos/:id/unpublish", controllers.UnpublishVideoPrompt)
			}
		}
	}
}
