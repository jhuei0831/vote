package model

import (
	"time"

	"github.com/google/uuid"
)

func (Invitation) TableName() string {
	return "invitations"
}

type Invitation struct {
	ID        	uint64    `gorm:"primary_key;auto_increment" json:"id"`
	SessionID   uuid.UUID `gorm:"index;not null;" json:"session_id"`
	CodeHash   	string    `gorm:"size:100;not null;" json:"code_hash"`
	Status    	bool      `gorm:"default:false;" json:"status"`
	CreatedAt 	time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Ballots   	[]Ballot  `gorm:"foreignKey:InvitationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballots,omitempty"`
}

type InvitationCreate struct {
	SessionID uuid.UUID `json:"session_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Number 		uint      `json:"number" binding:"required,min=1" example:"1"`
	Length 		uint      `json:"length" binding:"required,min=6" example:"8"`
	Format 		string    `json:"format" binding:"required,oneof=int en mix mixExcl mixLower mixUpper" example:"Aa1"`
}

type InvitationUpdate struct {
	SessionID   uuid.UUID `json:"session_id" example:"00000000-0000-0000-0000-000000000000"`
	Status    	bool      `json:"status" example:"false"`
}

type InvitationQuery struct {
	SessionID   uuid.UUID `json:"session_id" example:"00000000-0000-0000-0000-000000000000"`
	CodeHash 		string    `json:"code_hash" example:"code_hash"`
	Status   		bool      `json:"status" example:"false"`
	First    		int       `json:"first" binding:"min=0" example:"0"`
	After    		string    `json:"after" binding:"min=0" example:"0"`
	Last     		int       `json:"last" binding:"min=0" example:"0"`
	Before   		string    `json:"before" binding:"min=0" example:"0"`
}

type InvitationConnection struct {
	Edges      []InvitationEdge `json:"edges"`
	PageInfo 	 PageInfo   		`json:"pageInfo"`
	TotalCount int64			  	`json:"totalCount"`
}

type InvitationEdge struct {
	Node   Invitation   `json:"node"`
	Cursor string 			`json:"cursor"`
}

// GetFirst implements PaginationQuery
func (p *InvitationQuery) GetFirst() int {
	return p.First
}

// GetAfter implements PaginationQuery
func (p *InvitationQuery) GetAfter() string {
	return p.After
}

// GetLast implements PaginationQuery
func (p *InvitationQuery) GetLast() int {
	return p.Last
}

// GetBefore implements PaginationQuery
func (p *InvitationQuery) GetBefore() string {
	return p.Before
}

type VoterInfo struct {
	VoterID 		uint64
	SessionID  	uuid.UUID
	IsVoted 		bool
}

type VoterVerify struct {
	SessionID uuid.UUID `json:"session_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	CodeHash 	string    `json:"code_hash" binding:"required" example:"code_hash"`
}

type ValidateInviteResult struct {
	Success bool   `json:"success"`
	JWT     string `json:"jwt,omitempty"`
	Message string `json:"message,omitempty"`
}
