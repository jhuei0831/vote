package model

import (
	"time"
	"crypto/rand"
	"gorm.io/gorm"
	"math/big"
)

func (Vote) TableName() string {
	return "votes"
}

type Vote struct {
	ID    		uint64 	  `gorm:"primaryKey;autoIncrement:false;comment:ID" json:"id"`
	Title       string 	  `gorm:"size:100;not null;" json:"title"`
	Description string 	  `gorm:"size:255;not null;" json:"description"`
	StartTime   string 	  `gorm:"type:datetime;not null;" json:"start_time"`
	EndTime     string 	  `gorm:"type:datetime;not null;" json:"end_time"`
	UserID      uint64 	  `gorm:"index;not null;" json:"user_id"`
	CreatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
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

// This functions are called before creating any Post
func (v *Vote) BeforeCreate(tx *gorm.DB) (err error) {
	// 生成一個隨機的 ID
	v.ID, err = generateRandomID()
	return err
}

// generateRandomID 生成一个隨機的 uint64 ID
func generateRandomID() (uint64, error) {
	// 使用 crypto/rand 生成安全的隨機數
	n, err := rand.Int(rand.Reader, new(big.Int).SetUint64(^uint64(0)))
	if err != nil {
		return 0, err
	}

	return n.Uint64(), nil
}