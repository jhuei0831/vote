package model

import (
	"time"

	"github.com/google/uuid"
)

func (Password) TableName() string {
	return "passwords"
}

type Password struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	VoteID    uuid.UUID `gorm:"index;not null;" json:"vote_id"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Status    bool      `gorm:"default:false;" json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Ballots   []Ballot  `gorm:"foreignKey:PasswordID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballots,omitempty"`
}

type PasswordCreate struct {
	VoteID uuid.UUID `json:"vote_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Number uint      `json:"number" binding:"required,min=1" example:"1"`
	Length uint      `json:"length" binding:"required,min=6" example:"8"`
	Format string    `json:"format" binding:"required,oneof=int en mix mixExcl mixLower mixUpper" example:"Aa1"`
}

type PasswordUpdate struct {
	VoteID    	uuid.UUID `json:"voteId" example:"00000000-0000-0000-0000-000000000000"`
	Status    	bool      `json:"status" example:"false"`
}

type PasswordQuery struct {
	VoteID   uuid.UUID `json:"vote_id" example:"00000000-0000-0000-0000-000000000000"`
	Password string    `json:"password" example:"password"`
	Status   bool      `json:"status" example:"false"`
	First    int       `json:"first" binding:"min=0" example:"0"`
	After    string    `json:"after" binding:"min=0" example:"0"`
	Last     int       `json:"last" binding:"min=0" example:"0"`
	Before   string    `json:"before" binding:"min=0" example:"0"`
}

type PasswordConnection struct {
	Edges      []PasswordEdge `json:"edges"`
	PageInfo 	 PageInfo   		`json:"pageInfo"`
	TotalCount int64			  	`json:"totalCount"`
}

type PasswordEdge struct {
	Node   Password   `json:"node"`
	Cursor string 		`json:"cursor"`
}

// GetFirst implements PaginationQuery
func (p *PasswordQuery) GetFirst() int {
	return p.First
}

// GetAfter implements PaginationQuery
func (p *PasswordQuery) GetAfter() string {
	return p.After
}

// GetLast implements PaginationQuery
func (p *PasswordQuery) GetLast() int {
	return p.Last
}

// GetBefore implements PaginationQuery
func (p *PasswordQuery) GetBefore() string {
	return p.Before
}

type VoterInfo struct {
	VoterID uint64
	VoteID  uuid.UUID
	IsVoted bool
}

type VoterLogin struct {
	VoteID   uuid.UUID `json:"vote_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Password string    `json:"password" binding:"required" example:"password"`
}
