package repository

import (
	"vote/app/database"
	"vote/app/model"

	"github.com/google/uuid"
)

type PasswordRepository struct {
}

func NewPasswordRepository() PasswordRepository {
	return PasswordRepository{}
}

// GetVoteIdByVoterId 根據投票者ID獲取選票
func (p PasswordRepository) GetVoteIdByVoterId(voterId uint64) (uuid.UUID, error) {
	var voteId uuid.UUID
	err := database.SqlSession.Model(&model.Password{}).
		Select("vore_id").
		Where("id = ?", voterId).
		Scan(&voteId).Error

	if err != nil {
		return uuid.Nil, err
	}

	return voteId, nil
}

// GetPasswords 根據條件取得所有密碼。
func (p PasswordRepository) GetPasswords(passwordQuery *model.PasswordQuery) ([]model.Password, int64, error) {
	var passwords []model.Password
	var total int64

	query := database.SqlSession.Model(&model.Password{})

	// Status 精確查詢
	if passwordQuery.Status {
		query = query.Where("status = ?", passwordQuery.Status)
	}

	// 計算總筆數
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 取得資料
	// 使用分頁服務處理分頁
	paginationRepository := NewPaginationRepository[*model.PasswordQuery, model.Password]()
	query, err = paginationRepository.Handler(query, passwordQuery)
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&passwords).Error; err != nil {
		return nil, 0, err
	}

	// 查詢資料
	err = query.Find(&passwords).Error

	return passwords, total, err
}
