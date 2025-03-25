package model

import (
	"time"
)

func (Candidate) TableName() string {
	return "candidates"
}

type Candidate struct {
	ID 			uint64 		`gorm:"primary_key;auto_increment" json:"id"`
	QuestionID 	uint64 		`gorm:"index;not null;" json:"question_id"`
	Name 		string 		`gorm:"size:100;not null;" json:"name"`
	Result 		string 		`gorm:"default:null;" json:"result"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CandidateCreate struct {
	QuestionID uint64 `json:"question_id" binding:"required" example:"1"`
	Name       string `json:"name" binding:"required" example:"name"`
}