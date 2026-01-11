package models

import (
	"time"
)

type User struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Username          string     `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email             string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash      string     `gorm:"size:255;not null" json:"-"`
	FullName          string     `gorm:"size:255" json:"full_name"`
	Role              string     `gorm:"type:enum('user','admin');default:'user';not null" json:"role"`
	ProfilePictureURL string     `gorm:"size:500" json:"profile_picture_url"`
	Bio               string     `gorm:"type:text" json:"bio"`
	Interests         string     `gorm:"type:text" json:"interests"` // JSON array of interest IDs
	TotalCreations    int        `gorm:"default:0" json:"total_creations"`
	TotalLikes        int        `gorm:"default:0" json:"total_likes"`
	TrendingScore     int        `gorm:"default:0" json:"trending_score"`
	CommunityRank     *int       `json:"community_rank"`
	IsVerified        bool       `gorm:"default:false" json:"is_verified"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	EmailVerified     bool       `gorm:"default:false" json:"email_verified"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	LastLogin         *time.Time `json:"last_login"`
}

func (User) TableName() string {
	return "users"
}
