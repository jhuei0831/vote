package model

import (
	"time"

	"github.com/google/uuid"
)

func (Vote) TableName() string {
	return "votes"
}

type Vote struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Uuid        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();" json:"uuid"`
	Title       string     `gorm:"size:100;not null;" json:"title"`
	Description string     `gorm:"size:255;" json:"description"`
	StartTime   time.Time  `gorm:"not null;" json:"start_time"`
	EndTime     time.Time  `gorm:"not null;" json:"end_time"`
	UserID      uint64     `gorm:"index;not null;" json:"user_id"`
	Status      int        `gorm:"default:0;not null;" json:"status"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Questions   []Question `gorm:"foreignKey:VoteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"questions,omitempty"`
	Passwords   []Password `gorm:"foreignKey:VoteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"passwords,omitempty"`
}

type VoteCreate struct {
	Title       string    `json:"title" binding:"required,max=100" example:"title"`
	Description string    `json:"description" binding:"max=255" example:"description"`
	UserID      uint64    `json:"user_id" example:"1"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2006-01-02 15:04:05"`
	EndTime     time.Time `json:"end_time" binding:"required" example:"2006-01-02 15:04:05"`
}

type VoteUpdate struct {
	Title       string    `json:"title" binding:"required,max=100" example:"title"`
	Description string    `json:"description" binding:"max=255" example:"description"`
	UserID      uint64    `json:"user_id" example:"1"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2006-01-02 15:04:05"`
	EndTime     time.Time `json:"end_time" binding:"required" example:"2006-01-02 15:04:05"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Query parameters for filtering, sorting, and pagination
type VoteQuery struct {
	ID        uint64    `json:"id" example:"1"`
	Uuid      uuid.UUID `json:"uuid" example:"00000000-0000-0000-0000-000000000000"`
	Title     string    `json:"title" example:"title"`
	StartTime time.Time `json:"start_time" example:"2006-01-02 15:04:05"`
	EndTime   time.Time `json:"end_time" example:"2006-01-02 15:04:05"`
	First     int       `json:"first" binding:"min=1" example:"1"`
	After     string    `json:"after" binding:"min=1" example:"1"`
	Last      int       `json:"last" binding:"min=1" example:"1"`
	Before    string    `json:"before" binding:"min=1" example:"1"`
}

type VoteConnection struct {
	Edges      []VoteEdge `json:"edges"`
	PageInfo 	 PageInfo   `json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type VoteEdge struct {
	Node   Vote   `json:"node"`
	Cursor string `json:"cursor"`
}

// GetFirst implements PaginationQuery
func (q *VoteQuery) GetFirst() int {
	return q.First
}

// GetAfter implements PaginationQuery
func (q *VoteQuery) GetAfter() string {
	return q.After
}

// GetLast implements PaginationQuery
func (q *VoteQuery) GetLast() int {
	return q.Last
}

// GetBefore implements PaginationQuery
func (q *VoteQuery) GetBefore() string {
	return q.Before
}
