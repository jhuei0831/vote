package repository

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

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

	// VoteID 精確查詢
	if passwordQuery.VoteID != uuid.Nil {
		query = query.Where("vote_id = ?", passwordQuery.VoteID)
	}

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

// CreatePasswords 根據提供的投票ID和密碼列表創建密碼。
func (p PasswordRepository) CreatePasswords(voteId uuid.UUID, number uint, length uint, format string) ([]*model.Password, error) {
	passwordUtil := &utils.Password{}
	// Generate Passwords
	passwords, err := passwordUtil.GeneratePassword(number, length, format)
	if err != nil {
		return nil, err
	}

	// Encrypt Passwords
	passwordModels := make([]*model.Password, len(passwords))
	for i, password := range passwords {
		passwordEncrypt, err := passwordUtil.Encrypt(password)
		if err != nil {
			return nil, err
		}
		passwordModels[i] = &model.Password{
			VoteID:   voteId,
			Password: passwordEncrypt,
		}
	}

	// Use transaction to ensure all passwords are created successfully
	transaction := database.SqlSession.Begin()
	err = transaction.CreateInBatches(&passwordModels, 100).Error

	if err != nil {
		transaction.Rollback()
		return nil, err
	}

	return passwordModels, transaction.Commit().Error
}

// UpdatePasswordStatus 根據提供的投票ID和密碼ID列表更新密碼狀態。
func (p PasswordRepository) UpdatePasswordStatus(voteId uuid.UUID, passwordIDs []uint64, status bool) ([]*model.Password, error) {
	var passwordModels []*model.Password
	err := database.SqlSession.
		Model(&model.Password{}).
		Where("vote_id = ? AND id IN ?", voteId, passwordIDs).
		Update("status", status).
		Scan(&passwordModels).Error

	if err != nil {
		return nil, err
	}

	return passwordModels, nil
}

// DeletePassword 根據提供的密碼ID列表刪除密碼。
func (p PasswordRepository) DeletePassword(ids []uint64, userInfo model.UserInfo) ([]*model.Password, error) {
	var passwords []*model.Password

	query := database.SqlSession.
		Model(&model.Password{}).
		Where("passwords.id IN ?", ids)

	// 非管理員需檢查所屬 user
	if !userInfo.IsAdmin {
		query = query.
			Joins("JOIN votes ON passwords.vote_id = votes.uuid").
			Where("votes.user_id = ?", userInfo.UserID)
	}

	// 取得密碼並直接刪除
	if err := query.Find(&passwords).Error; err != nil {
		return nil, err
	}

	if len(passwords) == 0 {
		return passwords, nil
	}

	// 提取 IDs 並進行刪除
	passwordIds := make([]uint64, len(passwords))
	for i, password := range passwords {
		passwordIds[i] = password.ID
	}

	deleteError := database.SqlSession.Delete(&model.Password{}, passwordIds).Error
	return passwords, deleteError
}