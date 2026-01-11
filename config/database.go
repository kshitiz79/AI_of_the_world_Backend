package config

import (
	"ai-of-the-world-backend/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establishes connection to MySQL database
func ConnectDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ Failed to get database instance:", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate OTP table
	if err := DB.AutoMigrate(&models.OTP{}); err != nil {
		log.Println("⚠️  Failed to auto-migrate OTP table:", err)
	} else {
		log.Println("✅ OTP table migrated successfully")
	}

	log.Println("✅ Database connected successfully")
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("❌ Error getting database instance:", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Println("❌ Error closing database:", err)
		return
	}

	log.Println("✅ Database connection closed")
}
