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
	VoteID      uuid.UUID   `gorm:"index;not null;" json:"vote_id"`
	Title       string 		`gorm:"size:100;not null;" json:"title"`
	Description string 		`gorm:"size:255;not null;" json:"description"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type QuestionCreate struct {
	VoteID      uuid.UUID   `json:"vote_id" binding:"required" example:"1"`
	Title       string 		`json:"title" binding:"required" example:"title"`
	Description string 		`json:"description" binding:"required" example:"description"`
}