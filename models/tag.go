package models

import (
	"time"
)

type Tag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Category    string    `gorm:"type:enum('Style','Mood','Theme','Technique','Color','Other');default:'Other';not null" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}

// CreateTagRequest represents the request body for creating a tag
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Category    string `json:"category" binding:"required,oneof=Style Mood Theme Technique Color Other"`
	Description string `json:"description"`
}

// UpdateTagRequest represents the request body for updating a tag
type UpdateTagRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Category    string `json:"category" binding:"omitempty,oneof=Style Mood Theme Technique Color Other"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
