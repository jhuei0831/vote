package model

import (
	"github.com/google/uuid"
	"time"
)

func (Vote) TableName() string {
	return "votes"
}

type Vote struct {
	ID    		uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();comment:ID" json:"id"`
	Title       string 	  `gorm:"size:100;not null;" json:"title"`
	Description string 	  `gorm:"size:255;not null;" json:"description"`
	StartTime   string 	  `gorm:"not null;" json:"start_time"`
	EndTime     string 	  `gorm:"not null;" json:"end_time"`
	UserID      uint64 	  `gorm:"index;not null;" json:"user_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type VoteCreate struct {
	Title       string `json:"title" binding:"required" example:"title"`
	Description string `json:"description" binding:"required" example:"description"`
	UserID      uint64 `json:"user_id" example:"1"`
	StartTime   string `json:"startTime" binding:"required" example:"2006-01-02 15:04:05"`
	EndTime     string `json:"endTime" binding:"required" example:"2006-01-02 15:04:05"`
}

type VoteUpdate struct {
	Title       string `json:"title" example:"title"`
	Description string `json:"description" example:"description"`
	UserID      uint64 `json:"user_id" example:"1"`
	StartTime   string `json:"startTime" example:"2006-01-02 15:04:05"`
	EndTime     string `json:"endTime" example:"2006-01-02 15:04:05"`
}