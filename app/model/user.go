package model

import (
	"time"

	"gorm.io/gorm"
)

func (User) TableName() string {
	return "users"
}

type User struct {
	gorm.Model
	ID           uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Account      string    `gorm:"size:100;not null;unique" json:"account"`
	Password     string    `gorm:"size:100;not null;" json:"password"`
	Email        string    `gorm:"size:100;not null;unique" json:"email"`
	CreatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Votes        []Vote    `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}