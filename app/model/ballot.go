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
	PasswordID    uint64    	 		`gorm:"index;not null;" json:"password_id"`
	QuestionID	  uint64    	 		`gorm:"index;not null;" json:"question_id"`
	CreatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	  time.Time 	 		`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	BallotSelects []BallotSelect 	`gorm:"foreignKey:BallotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ballot_selects,omitempty"`
}

type BallotCandidate struct {
	CandidateID   uint64 `json:"candidateId" binding:"required"`
	IsSelected 		bool   `json:"isSelected" binding:"required"`
}

type BallotQuestions struct {
	QuestionID  uint64             `json:"questionId" binding:"required"`
	Candidates  []BallotCandidate  `json:"candidates"`
}

type BallotCreate struct {
	Selections []BallotQuestions `json:"selections" binding:"required"`
}

type BallotQuery struct {
	VoteID    	uuid.UUID   `json:"voteId" binding:"required"`
	QuestionID  uint64   		`json:"questionId"`
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