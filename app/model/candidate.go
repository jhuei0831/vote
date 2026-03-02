package model

import (
	"time"

	"github.com/google/uuid"
)

func (Candidate) TableName() string {
	return "candidates"
}

type Candidate struct {
	ID 					uint64 			`gorm:"primary_key;auto_increment" json:"id"`
	QuestionID 	uint64 			`gorm:"index;not null;" json:"question_id"`
	Name 				string 			`gorm:"size:100;not null;" json:"name"`
	Result 			string 			`gorm:"default:null;" json:"result"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CandidateCreate struct {
	QuestionID 	uint64 			`json:"question_id" binding:"required" example:"1"`
	Name       	string 			`json:"name" binding:"required,max=100" example:"name"`
}

type CandidateUpdate struct {
	QuestionID 	uint64 			`json:"question_id" binding:"required" example:"1"`
	Name 				string 			`json:"name" binding:"max=100" example:"name"`
}

type CandidateQuery struct {	
	VoteID 			uuid.UUID 	`json:"vote_id" example:"00000000-0000-0000-0000-000000000000"`
	QuestionID 	uint64 			`json:"question_id" example:"1"`
	Name	   		string 			`json:"name" example:"name"`
	First     	int       	`json:"first" binding:"min=1" example:"1"`
	After     	string    	`json:"after" binding:"min=1" example:"1"`
	Last      	int       	`json:"last" binding:"min=1" example:"1"`
	Before    	string    	`json:"before" binding:"min=1" example:"1"`
}

type CandidateConnection struct {
	Edges      []CandidateEdge `json:"edges"`
	PageInfo 	 PageInfo   `json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type CandidateEdge struct {
	Node   Candidate   `json:"node"`
	Cursor string `json:"cursor"`
}

// GetFirst implements PaginationQuery
func (c *CandidateQuery) GetFirst() int {
	return c.First
}

// GetAfter implements PaginationQuery
func (c *CandidateQuery) GetAfter() string {
	return c.After
}

// GetLast implements PaginationQuery
func (c *CandidateQuery) GetLast() int {
	return c.Last
}

// GetBefore implements PaginationQuery
func (c *CandidateQuery) GetBefore() string {
	return c.Before
}