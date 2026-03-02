package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"
	"vote/app/utils"

	"github.com/google/uuid"
)

type PasswordService struct {
}

func NewPasswordService() PasswordService {
	return PasswordService{}
}

// GetPassword 根據提供的投票ID和密碼，檢查密碼是否存在。
func (p PasswordService) GetPassword(voteId uuid.UUID, password string) (*model.Password, error) {
	passwordModel := model.Password{}
	err := database.SqlSession.
		Where("vote_id = ? AND password = ? AND status = true", voteId, password).
		First(&passwordModel).
		Error

	if err != nil {
		return nil, err
	}

	return &passwordModel, nil
}

// GetPassword 根據提供的投票ID，檢索所有密碼。
func (p PasswordService) GetPasswords(voteId uuid.UUID, passwordQuery *model.PasswordQuery) ([]*model.PasswordConnection, error) {
	passwords, total, err := repository.NewPasswordRepository().GetPasswords(passwordQuery)
	if err != nil {
		return nil, err
	}

	paginationRepository := repository.NewPaginationRepository[*model.PasswordQuery, model.Password]()
	passwords, hasPreviousPage, hasNextPage := paginationRepository.HasPreviousNextPage(passwords, passwordQuery)

	paginationService := NewPaginationService[model.Password, model.PasswordEdge, *model.PasswordConnection]()
	connection := paginationService.BuildConnection(passwords, total, hasPreviousPage, hasNextPage,
		func(password model.Password) uint64 {
			return password.ID
		},
	)

	return []*model.PasswordConnection{connection}, nil
}

// CreatePassword Create passwords can encrypt and decrypt
func (p PasswordService) CreatePassword(voteId uuid.UUID, number uint, length uint, format string) ([]*model.Password, error) {
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

// UpdatePasswordStatus 更新密碼狀態
func (p PasswordService) UpdatePasswordStatus(voteId uuid.UUID, passwordIDs []uint64, status bool) ([]*model.Password, error) {
	var passwordModels []*model.Password
	err := database.SqlSession.
		Where("vote_id = ? AND id IN ?", voteId, passwordIDs).
		Update("status", status).
		Scan(&passwordModels).Error

	if err != nil {
		return nil, err
	}

	return passwordModels, nil
}

// DeletePassword 根據提供的密碼ID列表刪除密碼。
func (p PasswordService) DeletePassword(ids []uint64, userInfo model.UserInfo) ([]*model.Password, error) {
	var passwordModels []*model.Password
	err := database.SqlSession.
		Where("id IN ?", ids).
		Delete(&passwordModels).Error

	if err != nil {
		return nil, err
	}

	return passwordModels, nil
}
