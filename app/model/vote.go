package model

import (
	"github.com/google/uuid"
	"time"
)

func (Vote) TableName() string {
	return "votes"
}

type Vote struct {
	ID    			uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();comment:ID" json:"id"`
	Title       string 	   `gorm:"size:100;not null;" json:"title"`
	Description string 	   `gorm:"size:255;" json:"description"`
	StartTime   time.Time  `gorm:"not null;" json:"start_time"`
	EndTime     time.Time  `gorm:"not null;" json:"end_time"`
	UserID      uint64 	   `gorm:"index;not null;" json:"user_id"`
	Status			int	   	   `gorm:"default:0;not null;" json:"status"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Questions   []Question `gorm:"foreignKey:VoteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"questions,omitempty"`
	Passwords   []Password `gorm:"foreignKey:VoteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"passwords,omitempty"`
}

type VoteCreate struct {
	Title       string 	  `json:"title" binding:"required,max=100" example:"title"`
	Description string 	  `json:"description" binding:"max=255" example:"description"`
	UserID      uint64 	  `json:"user_id" example:"1"`
	StartTime   time.Time `json:"startTime" binding:"required" example:"2006-01-02 15:04:05"`
	EndTime     time.Time `json:"endTime" binding:"required" example:"2006-01-02 15:04:05"`
}

type VoteUpdate struct {
	Title       string 	  `json:"title" binding:"required,max=100" example:"title"`
	Description string 	  `json:"description" binding:"max=255" example:"description"`
	UserID      uint64 	  `json:"user_id" example:"1"`
	StartTime   time.Time `json:"startTime" binding:"required" example:"2006-01-02 15:04:05"`
	EndTime     time.Time `json:"endTime" binding:"required" example:"2006-01-02 15:04:05"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Query parameters for filtering, sorting, and pagination
type VoteQuery struct {
	ID	  			uuid.UUID 	`json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Title	  		string    	`json:"title" example:"title"`
	StartTime 	time.Time 	`json:"start_time" example:"2006-01-02 15:04:05"`
	EndTime   	time.Time 	`json:"end_time" example:"2006-01-02 15:04:05"`
	Page	 			int    			`form:"page,default=1" json:"page" binding:"min=1" example:"1"`
	Size	 			int    			`form:"size,default=1" json:"size" binding:"min=1" example:"10"`
}