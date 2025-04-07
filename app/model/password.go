package model

import (
	"time"

	"github.com/google/uuid"
)

func (Password) TableName() string {
	return "passwords"
}

type Password struct {
	ID        uint64       `gorm:"primary_key;auto_increment" json:"id"`
	VoteID	  uuid.UUID    `gorm:"index;not null;" json:"vote_id"`
	Password  string       `gorm:"size:100;not null;" json:"password"`
	Status	  bool         `gorm:"default:false;" json:"status"`
	CreatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Ballots	  []Ballot     `gorm:"foreignKey:PasswordID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballots,omitempty"`
}

type PasswordCreate struct {
	VoteID uuid.UUID `json:"vote_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Number int       `json:"number" binding:"required,min=1" example:"1"`
	Length int       `json:"length" binding:"required,min=6" example:"8"`
	Format string    `json:"format" binding:"required,oneof=int en mix mixExcl mixLower mixUpper" example:"Aa1"`
}

type AnonLogin struct {
	VoteID   uuid.UUID `json:"vote_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Password string    `json:"password" binding:"required" example:"password"`
}

type PasswordQuery struct {
	VoteID 		uuid.UUID 	`json:"vote_id" example:"00000000-0000-0000-0000-000000000000"`
	Password 	string    	`json:"password" example:"password"`
	Status 		bool      	`json:"status" example:"false"`
	Page	 	int    		`form:"page,default=1" json:"page" binding:"min=1" example:"1"`
	Size	 	int    		`form:"size,default=10" json:"size" binding:"min=1" example:"10"`
}