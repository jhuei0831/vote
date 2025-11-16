package model

import (
	"time"
)

func (Ballot) TableName() string {
	return "ballots"
}

type Ballot struct {
	ID        	  uint64    	 		`gorm:"primary_key;auto_increment" json:"id"`
	PasswordID    uint64    	 		`gorm:"index;not null;" json:"password_id"`
	QuestionID	  uint64    	 		`gorm:"index;not null;" json:"question_id"`
	CreatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	BallotSelects []BallotSelect 	`gorm:"foreignKey:BallotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballot_selects,omitempty"`
}

type BallotCreate struct {
	Selections map[string][]string `json:"selections" binding:"required"`
}
