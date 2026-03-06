package model

func (BallotSelect) TableName() string {
	return "ballot_selects"
}

type BallotSelect struct {
	ID        	  uint64    	`gorm:"primary_key;auto_increment" json:"id"`
	BallotID      uint64    	`gorm:"index;not null;" json:"ballot_id"`
	PollOptionID	uint64    	`gorm:"index;not null;" json:"poll_option_id"`
}