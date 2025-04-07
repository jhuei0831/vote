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
func (v VoteService) SelectAllVotes(isAdmin bool, userId uint64, voteQuery model.VoteQuery) ([]model.Vote, int64, error) {
	var votes []model.Vote
	var total int64
	query := database.SqlSession.Model(&model.Vote{})

	if !isAdmin {
		query = query.Where("user_id = ?", userId)
	}

	// 設定查詢條件
	page := voteQuery.Page
	size := voteQuery.Size

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果有 page 和 size，加入分頁條件
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	// 查詢資料
	err = query.Find(&votes).Error
	return votes, total, err
}

// CreateOneVote 創建新的投票。
func (v VoteService) CreateVote(form model.VoteCreate) (model.Vote, error) {
	vote := model.Vote{
		Title:       form.Title,
		Description: form.Description,
		UserID:      form.UserID,
		StartTime:   form.StartTime,
		EndTime:     form.EndTime,
	}

	insertErr := database.SqlSession.Model(&model.Vote{}).Create(&vote).Error
	return vote, insertErr
}

// UpdateOneVote 更新投票。
func (v VoteService) UpdateVote(id uuid.UUID, form model.VoteUpdate) (model.Vote, error) {
	var vote model.Vote
	// 更新投票
	update := database.SqlSession.Model(&vote).
		Clauses(clause.Returning{}).
		Where("id=?", id).
		Updates(&form)
	
	return vote, update.Error
}

// DeleteOneVote 刪除投票。
func (v VoteService) DeleteVote(voteIds []uuid.UUID, isAdmin bool, userId uint64) ([]model.Vote, error) {
	var votes []model.Vote
	query := database.SqlSession.
		Clauses(clause.Returning{}).
		Where("id IN (?)", voteIds)

	if !isAdmin {
		query = query.Where("user_id = ?", userId)
	}

	err := query.Delete(&votes).Error
	return votes, err
}