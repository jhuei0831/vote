package service

import (
	"vote/app/database"
	"vote/app/model"
)

type VoteService struct {
}

func NewVoteService() VoteService {
	return VoteService{}
}

// SelectOneVote 根據提供的 ID 檢查投票是否存在。
func (v VoteService) SelectOneVote(id int64) (*model.Vote, error) {
	voteOne := &model.Vote{}
	err := database.SqlSession.
		Select([]string{"id", "title", "description", "start_time", "end_time"}).
		Where("id=?", id).
		First(&voteOne).Error
	if err != nil {
		return nil, err
	} else {
		return voteOne, nil
	}
}

// SelectAllVotes 檢索所有投票。
func (v VoteService) SelectAllVotes(isAdmin bool, user_id uint64) ([]model.Vote, error) {
	var votes []model.Vote
	var err error
	if isAdmin {
		err = database.SqlSession.Find(&votes).Error
	} else {
		err = database.SqlSession.Where("user_id = ?", user_id).Find(&votes).Error
	}

	return votes, err
}

// CreateOneVote 創建一個新的投票。
func (v VoteService) CreateOneVote(title string, description string, userId uint64, startTime string, endTime string) error {
	vote := model.Vote{
		Title:       title,
		Description: description,
		UserID:      userId,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	insertErr := database.SqlSession.Model(&model.Vote{}).Create(&vote).Error
	return insertErr
}