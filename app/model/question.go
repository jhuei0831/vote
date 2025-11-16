package model

import (
	"time"

	"github.com/google/uuid"
)

func (Question) TableName() string {
	return "questions"
}

type Question struct {
	ID          uint64 			`gorm:"primary_key;auto_increment" json:"id"`
	VoteID      uuid.UUID   `gorm:"index;type:uuid;not null;" json:"vote_id"`
	Title       string 			`gorm:"size:100;not null;" json:"title"`
	Description string 			`gorm:"size:255;" json:"description"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Candidates  []Candidate `gorm:"foreignKey:QuestionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"candidates,omitempty"`
}

type QuestionCreate struct {
	VoteID      uuid.UUID   `json:"vote_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Title       string 			`json:"title" binding:"required" example:"title"`
	Description string 			`json:"description" example:"description"`
}

// Query parameters for filtering, sorting, and pagination
type QuestionQuery struct {
	VoteID  		uuid.UUID 	`json:"vote_id" example:"00000000-0000-0000-0000-000000000000"`
	Title	  		string    	`json:"title" example:"title"`
	Candidates  bool 				`form:"candidates,default=false" json:"candidates" example:"false"`
	First     	int       	`json:"first" binding:"min=1" example:"1"`
	After     	string    	`json:"after" binding:"min=1" example:"1"`
	Last      	int       	`json:"last" binding:"min=1" example:"1"`
	Before    	string    	`json:"before" binding:"min=1" example:"1"`
}

type QuestionConnection struct {
	Edges      []QuestionEdge `json:"edges"`
	PageInfo 	 PageInfo   `json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type QuestionEdge struct {
	Node   Question   `json:"node"`
	Cursor string `json:"cursor"`
}

// GetFirst implements PaginationQuery
func (q *QuestionQuery) GetFirst() int {
	return q.First
}

// GetAfter implements PaginationQuery
func (q *QuestionQuery) GetAfter() string {
	return q.After
}

// GetLast implements PaginationQuery
func (q *QuestionQuery) GetLast() int {
	return q.Last
}

// GetBefore implements PaginationQuery
func (q *QuestionQuery) GetBefore() string {
	return q.Before
}
