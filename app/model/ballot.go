package model

import (
	"time"

	"github.com/google/uuid"
)

func (Ballot) TableName() string {
	return "ballots"
}

type Ballot struct {
	ID        	  uint64    	 		`gorm:"primary_key;auto_increment" json:"id"`
	InvitationID  uint64    	 		`gorm:"index;not null;" json:"invitation_id"`
	PollID	  		uint64    	 		`gorm:"index;not null;" json:"poll_id"`
	CreatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	BallotSelects []BallotSelect 	`gorm:"foreignKey:BallotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballot_selects,omitempty"`
}

type BallotPollOptions struct {
	PollOptionID   	uint64 `json:"pollOptionId" binding:"required"`
	IsSelected 			bool   `json:"isSelected" binding:"required"`
}

type BallotPolls struct {
	PollID  		 uint64             `json:"pollId" binding:"required"`
	PollOptions  []BallotPollOptions  `json:"pollOptions"`
}

type BallotCreate struct {
	Selections []BallotPolls `json:"selections" binding:"required"`
}

type BallotQuery struct {
	SessionID   uuid.UUID   `json:"sessionId" binding:"required"`
	PollID  		uint64   		`json:"pollId"`
	VoterID    	uint64   		`json:"voterId"`
	First     	int       	`json:"first" binding:"min=1" example:"1"`
	After     	string    	`json:"after" binding:"min=1" example:"1"`
	Last      	int       	`json:"last" binding:"min=1" example:"1"`
	Before    	string    	`json:"before" binding:"min=1" example:"1"`
}

type BallotConnection struct {
	Edges      []BallotEdge `json:"edges"`
	PageInfo 	 PageInfo   	`json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type BallotEdge struct {
	Node   Ballot   `json:"node"`
	Cursor string 	`json:"cursor"`
}

// GetFirst implements PaginationQuery
func (b *BallotQuery) GetFirst() int {
	return b.First
}

// GetAfter implements PaginationQuery
func (b *BallotQuery) GetAfter() string {
	return b.After
}

// GetLast implements PaginationQuery
func (b *BallotQuery) GetLast() int {
	return b.Last
}

// GetBefore implements PaginationQuery
func (b *BallotQuery) GetBefore() string {
	return b.Before
}