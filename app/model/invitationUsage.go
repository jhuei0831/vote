package model

import (
	"time"

	"github.com/google/uuid"
)

// +gorm:"uniqueIndex:idx_invitation_vote"
type InvitationUsage struct {
	ID            uint64 				`gorm:"primaryKey" json:"id"`
	InvitationID  uint64 				`gorm:"uniqueIndex:idx_invitation_vote" json:"invitation_id"`
	VoterTempID   uuid.UUID 		`gorm:"type:uuid;uniqueIndex;uniqueIndex:idx_invitation_vote" json:"voter_temp_id"`
	UsedAt        time.Time 		`gorm:"default:CURRENT_TIMESTAMP" json:"used_at"`
	VoteCompleted bool   				`gorm:"default:false" json:"vote_completed"`
}
