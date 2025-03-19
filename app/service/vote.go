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
func (v VoteService) SelectAllVotes(isAdmin bool, userId uint64) ([]model.Vote, error) {
	var votes []model.Vote
	var err error
	if isAdmin {
		err = database.SqlSession.Find(&votes).Error
	} else {
		err = database.SqlSession.Where("user_id = ?", userId).Find(&votes).Error
	}

	return votes, err
}

// CreateOneVote 創建一個新的投票。
func (v VoteService) CreateOneVote(form model.VoteCreate) error {
	vote := model.Vote{
		Title:       form.Title,
		Description: form.Description,
		UserID:      form.UserID,
		StartTime:   form.StartTime,
		EndTime:     form.EndTime,
	}

	insertErr := database.SqlSession.Model(&model.Vote{}).Create(&vote).Error
	return insertErr
}

// UpdateOneVote 更新一個投票。
func (v VoteService) UpdateOneVote(id uint64, form model.VoteUpdate) error {
	updateErr := database.SqlSession.Model(&model.Vote{}).Where("id=?", id).Updates(&form).Error
	return updateErr
}