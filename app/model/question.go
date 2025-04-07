package model

import (
	"time"

	"github.com/google/uuid"
)

func (Question) TableName() string {
	return "questions"
}

type Question struct {
	ID          uint64 		`gorm:"primary_key;auto_increment" json:"id"`
	VoteID      uuid.UUID   `gorm:"index;type:uuid;not null;" json:"vote_id"`
	Title       string 		`gorm:"size:100;not null;" json:"title"`
	Description string 		`gorm:"size:255;" json:"description"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Candidates  []Candidate `gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"candidates,omitempty"`
}

type QuestionCreate struct {
	VoteID      uuid.UUID   `json:"vote_id" binding:"required" example:"1"`
	Title       string 		`json:"title" binding:"required" example:"title"`
	Description string 		`json:"description" example:"description"`
}

// Query parameters for filtering, sorting, and pagination
type QuestionQuery struct {
	Title	  	string    	`json:"title" example:"title"`
	Page	 	int    		`form:"page,default=1" json:"page" binding:"min=1" example:"1"`
	Size	 	int    		`form:"size,default=1" json:"size" binding:"min=1" example:"10"`
	Candidates  bool 		`form:"candidates,default=true" json:"candidates" example:"true"`
}