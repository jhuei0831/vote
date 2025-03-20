package service

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type VoteService struct {
}

func NewVoteService() VoteService {
	return VoteService{}
}

// SelectOneVote 根據提供的 ID 檢查投票是否存在。
func (v VoteService) SelectOneVote(id uuid.UUID) (*model.Vote, error) {
	voteOne := &model.Vote{}
	err := database.SqlSession.
		Select([]string{"id", "title", "description", "user_id", "start_time", "end_time"}).
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
func (v VoteService) CreateVote(form model.VoteCreate) error {
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
func (v VoteService) UpdateVote(id uuid.UUID, form model.VoteUpdate) (model.Vote, error) {
	var vote model.Vote
	// 更新投票
	update := database.SqlSession.Model(&vote).
		Clauses(clause.Returning{}).
		Where("id=?", id).
		Updates(&form)
	
	return vote, update.Error
}

// DeleteOneVote 刪除一個投票。
func (v VoteService) DeleteVote(voteIds []uuid.UUID, isAdmin bool, userId uint64) ([]model.Vote, error) {
	var votes []model.Vote
	var err error
	if isAdmin {
		err = database.SqlSession.
			Clauses(clause.Returning{}).
			Where("id IN (?)", voteIds).
			Delete(&votes).
			Error

		return votes, err
	}

	err = database.SqlSession.
		Clauses(clause.Returning{}).
		Where("id IN (?)", voteIds).
		Where("user_id = ?", userId).
		Delete(&votes).
		Error

	return votes, err
}