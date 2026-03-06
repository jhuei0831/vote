package model

import (
	"time"

	"github.com/google/uuid"
)

func (Poll) TableName() string {
	return "polls"
}

type Poll struct {
	ID          	uint64 				`gorm:"primary_key;auto_increment" json:"id"`
	SessionID   	uuid.UUID   	`gorm:"index;type:uuid;not null;" json:"session_id"`
	Title       	string 				`gorm:"size:100;not null;" json:"title"`
	Description 	string 				`gorm:"size:255;" json:"description"`
	CreatedAt   	time.Time   	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   	time.Time   	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	PollOptions  	[]PollOption 	`gorm:"foreignKey:PollID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"poll_options,omitempty"`
}

type PollCreate struct {
	SessionID   uuid.UUID   `json:"session_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Title       string 			`json:"title" binding:"required,max=100" example:"title"`
	Description string 			`json:"description" example:"description"`
}

type PollUpdate struct {
	SessionID   uuid.UUID   `json:"session_id" binding:"required" example:"00000000-0000-0000-0000-000000000000"`
	Title       string 			`json:"title" binding:"required,max=100" example:"title"`
	Description string 			`json:"description" example:"description"`
}

// Query parameters for filtering, sorting, and pagination
type PollQuery struct {
	SessionID  		uuid.UUID 	`json:"session_id" example:"00000000-0000-0000-0000-000000000000"`
	Title	  			string    	`json:"title" example:"title"`
	PollOptions  	bool 				`form:"poll_options,default=false" json:"poll_options" example:"false"`
	First     		int       	`json:"first" binding:"min=1" example:"1"`
	After     		string    	`json:"after" binding:"min=1" example:"1"`
	Last      		int       	`json:"last" binding:"min=1" example:"1"`
	Before    		string    	`json:"before" binding:"min=1" example:"1"`
}

type PollConnection struct {
	Edges      []PollEdge 	`json:"edges"`
	PageInfo 	 PageInfo   	`json:"pageInfo"`
	TotalCount int64			  `json:"totalCount"`
}

type PollEdge struct {
	Node   Poll   `json:"node"`
	Cursor string `json:"cursor"`
}

type PollList struct {
	List []Poll `json:"list"`
}

// GetFirst implements PaginationQuery
func (q *PollQuery) GetFirst() int {
	return q.First
}

// GetAfter implements PaginationQuery
func (q *PollQuery) GetAfter() string {
	return q.After
}

// GetLast implements PaginationQuery
func (q *PollQuery) GetLast() int {
	return q.Last
}

// GetBefore implements PaginationQuery
func (q *PollQuery) GetBefore() string {
	return q.Before
}
