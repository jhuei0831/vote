package model

import (
	"time"

	"github.com/google/uuid"
)

func (PollOption) TableName() string {
	return "poll_options"
}

type PollOption struct {
	ID 					uint64 			`gorm:"primary_key;auto_increment" json:"id"`
	PollID 			uint64 			`gorm:"index;not null;" json:"poll_id"`
	Name 				string 			`gorm:"size:100;not null;" json:"name"`
	Result 			string 			`gorm:"default:null;" json:"result"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type PollOptionCreate struct {
	PollID 			uint64 			`json:"poll_id" binding:"required" example:"1"`
	Name       	string 			`json:"name" binding:"required,max=100" example:"name"`
}

type PollOptionUpdate struct {
	PollID 			uint64 			`json:"poll_id" binding:"required" example:"1"`
	Name 				string 			`json:"name" binding:"max=100" example:"name"`
}

type PollOptionQuery struct {	
	SessionID 	uuid.UUID 	`json:"session_id" example:"00000000-0000-0000-0000-000000000000"`
	PollID 			uint64 			`json:"poll_id" example:"1"`
	Name	   		string 			`json:"name" example:"name"`
	First     	int       	`json:"first" binding:"min=1" example:"1"`
	After     	string    	`json:"after" binding:"min=1" example:"1"`
	Last      	int       	`json:"last" binding:"min=1" example:"1"`
	Before    	string    	`json:"before" binding:"min=1" example:"1"`
}

type PollOptionConnection struct {
	Edges      []PollOptionEdge `json:"edges"`
	PageInfo 	 PageInfo   `json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type PollOptionEdge struct {
	Node   PollOption   `json:"node"`
	Cursor string `json:"cursor"`
}

// GetFirst implements PaginationQuery
func (c *PollOptionQuery) GetFirst() int {
	return c.First
}

// GetAfter implements PaginationQuery
func (c *PollOptionQuery) GetAfter() string {
	return c.After
}

// GetLast implements PaginationQuery
func (c *PollOptionQuery) GetLast() int {
	return c.Last
}

// GetBefore implements PaginationQuery
func (c *PollOptionQuery) GetBefore() string {
	return c.Before
}