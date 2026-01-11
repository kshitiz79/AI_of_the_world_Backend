package models

import (
	"time"
)

// ImagePrompt represents an image submission
type ImagePrompt struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          uint       `gorm:"not null" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectTitle    string     `gorm:"size:100;not null" json:"project_title"`
	Prompt          string     `gorm:"type:text;not null" json:"prompt"`
	TechnicalNotes  string     `gorm:"type:text" json:"technical_notes"`
	ModelOrTool     string     `gorm:"size:255" json:"model_or_tool"`
	CreatorCredit   string     `gorm:"size:255;not null" json:"creator_credit"`
	ImageURL        string     `gorm:"size:500;not null" json:"image_url"`
	ImageFilename   string     `gorm:"size:255" json:"image_filename"`
	ImageSizeBytes  *int       `json:"image_size_bytes"`
	ImageWidth      *int       `json:"image_width"`
	ImageHeight     *int       `json:"image_height"`
	Status          string     `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"status"`
	VerifiedBy      *uint      `json:"verified_by"`
	VerifiedAt      *time.Time `json:"verified_at"`
	RejectionReason string     `gorm:"type:text" json:"rejection_reason"`
	LikesCount      int        `gorm:"default:0" json:"likes_count"`
	ViewsCount      int        `gorm:"default:0" json:"views_count"`
	DownloadsCount  int        `gorm:"default:0" json:"downloads_count"`
	IsFeatured      bool       `gorm:"default:false" json:"is_featured"`
	IsPublished     bool       `gorm:"default:false" json:"is_published"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Tags            []Tag      `gorm:"many2many:image_prompt_tags;" json:"tags,omitempty"`
}

func (ImagePrompt) TableName() string {
	return "image_prompts"
}

// GIFPrompt represents a GIF submission
type GIFPrompt struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserID             uint       `gorm:"not null" json:"user_id"`
	User               User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectTitle       string     `gorm:"size:100;not null" json:"project_title"`
	Prompt             string     `gorm:"type:text;not null" json:"prompt"`
	TechnicalNotes     string     `gorm:"type:text" json:"technical_notes"`
	ModelOrTool        string     `gorm:"size:255" json:"model_or_tool"`
	CreatorCredit      string     `gorm:"size:255;not null" json:"creator_credit"`
	GIFURL             string     `gorm:"size:500;not null;column:gif_url" json:"gif_url"`
	GIFFilename        string     `gorm:"size:255;column:gif_filename" json:"gif_filename"`
	GIFSizeBytes       *int       `gorm:"column:gif_size_bytes" json:"gif_size_bytes"`
	GIFWidth           *int       `gorm:"column:gif_width" json:"gif_width"`
	GIFHeight          *int       `gorm:"column:gif_height" json:"gif_height"`
	GIFDurationSeconds *float64   `gorm:"column:gif_duration_seconds" json:"gif_duration_seconds"`
	GIFFrameCount      *int       `gorm:"column:gif_frame_count" json:"gif_frame_count"`
	Status             string     `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"status"`
	VerifiedBy         *uint      `json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	RejectionReason    string     `gorm:"type:text" json:"rejection_reason"`
	LikesCount         int        `gorm:"default:0" json:"likes_count"`
	ViewsCount         int        `gorm:"default:0" json:"views_count"`
	DownloadsCount     int        `gorm:"default:0" json:"downloads_count"`
	IsFeatured         bool       `gorm:"default:false" json:"is_featured"`
	IsPublished        bool       `gorm:"default:false" json:"is_published"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Tags               []Tag      `gorm:"many2many:gif_prompt_tags;" json:"tags,omitempty"`
}

func (GIFPrompt) TableName() string {
	return "gif_prompts"
}

// VideoPrompt represents a video submission
type VideoPrompt struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	UserID               uint       `gorm:"not null" json:"user_id"`
	User                 User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectTitle         string     `gorm:"size:100;not null" json:"project_title"`
	Prompt               string     `gorm:"type:text;not null" json:"prompt"`
	TechnicalNotes       string     `gorm:"type:text" json:"technical_notes"`
	ModelOrTool          string     `gorm:"size:255" json:"model_or_tool"`
	CreatorCredit        string     `gorm:"size:255;not null" json:"creator_credit"`
	VideoURL             string     `gorm:"size:500;not null;column:video_url" json:"video_url"`
	VideoFilename        string     `gorm:"size:255;column:video_filename" json:"video_filename"`
	VideoSizeBytes       *int       `gorm:"column:video_size_bytes" json:"video_size_bytes"`
	VideoWidth           *int       `gorm:"column:video_width" json:"video_width"`
	VideoHeight          *int       `gorm:"column:video_height" json:"video_height"`
	VideoDurationSeconds *float64   `gorm:"column:video_duration_seconds" json:"video_duration_seconds"`
	VideoFormat          string     `gorm:"size:50;column:video_format" json:"video_format"`
	VideoFPS             *int       `gorm:"column:video_fps" json:"video_fps"`
	Status               string     `gorm:"type:enum('pending','approved','rejected');default:'pending';not null" json:"status"`
	VerifiedBy           *uint      `json:"verified_by"`
	VerifiedAt           *time.Time `json:"verified_at"`
	RejectionReason      string     `gorm:"type:text" json:"rejection_reason"`
	LikesCount           int        `gorm:"default:0" json:"likes_count"`
	ViewsCount           int        `gorm:"default:0" json:"views_count"`
	DownloadsCount       int        `gorm:"default:0" json:"downloads_count"`
	IsFeatured           bool       `gorm:"default:false" json:"is_featured"`
	IsPublished          bool       `gorm:"default:false" json:"is_published"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Tags                 []Tag      `gorm:"many2many:video_prompt_tags;" json:"tags,omitempty"`
}

func (VideoPrompt) TableName() string {
	return "video_prompts"
}
