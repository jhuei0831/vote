package repository

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type VoteRepository struct {
}

func NewVoteRepository() VoteRepository {
	return VoteRepository{}
}

// GetVoteByUUID 根據提供的UUID取得投票。
func (v VoteRepository) GetVoteByUUID(uuid uuid.UUID) (*model.Vote, error) {
	vote := &model.Vote{}

	err := database.SqlSession.
		Where("uuid=?", uuid).
		Preload("Questions").
		First(&vote).Error

	return vote, err
}

// GetVotes 根據條件取得所有投票。
func (v VoteRepository) GetVotes(userInfo model.UserInfo, voteQuery *model.VoteQuery) ([]model.Vote, int64, error) {
	var votes []model.Vote
	var total int64

	query := database.SqlSession.Model(&model.Vote{}).Preload("Questions")

	if !userInfo.IsAdmin {
		query = query.Where("user_id = ?", userInfo.UserID)
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 使用分頁服務處理分頁
	paginationRepository := NewPaginationRepository[*model.VoteQuery, model.Vote]()
	query, err = paginationRepository.Handler(query, voteQuery)
	if err != nil {
		return nil, 0, err
	}

	// 查詢資料
	err = query.Find(&votes).Error

	return votes, total, err
}

// CreateVote 創建新的投票。
func (v VoteRepository) CreateVote(form model.VoteCreate) (*model.Vote, error) {
	vote := model.Vote{
		Title:       form.Title,
		Description: form.Description,
		UserID:      form.UserID,
		StartTime:   form.StartTime,
		EndTime:     form.EndTime,
	}

	insertErr := database.SqlSession.Create(&vote).Error

	return &vote, insertErr
}

// UpdateVote 更新現有的投票。
func (v VoteRepository) UpdateVote(uuid uuid.UUID, form model.VoteUpdate) (*model.Vote, error) {
	var vote model.Vote

	updateError := database.SqlSession.Model(&vote).
		Clauses(clause.Returning{}).
		Where("uuid=?", uuid).
		Updates(&form).Error

	return &vote, updateError
}

// DeleteVotes 刪除投票。
func (v VoteRepository) DeleteVotes(voteUuids []uuid.UUID) ([]*model.Vote, error) {
	var votes []*model.Vote

	deleteErr := database.SqlSession.
		Clauses(clause.Returning{}).
		Where("uuid IN ?", voteUuids).
		Delete(&votes).Error

	return votes, deleteErr
}
