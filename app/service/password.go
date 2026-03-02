package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/repository"

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
	passwords, err := repository.NewPasswordRepository().CreatePasswords(voteId, number, length, format)
	if err != nil {
		return nil, err
	}

	return passwords, nil
}

// UpdatePasswordStatus 更新密碼狀態
func (p PasswordService) UpdatePasswordStatus(voteId uuid.UUID, passwordIDs []uint64, status bool) ([]*model.Password, error) {
	passwords, err := repository.NewPasswordRepository().UpdatePasswordStatus(voteId, passwordIDs, status)
	if err != nil {
		return nil, err
	}

	return passwords, nil
}

// DeletePassword 根據提供的密碼ID列表刪除密碼。
func (p PasswordService) DeletePassword(ids []uint64, userInfo model.UserInfo) ([]*model.Password, error) {
	passwords, err := repository.NewPasswordRepository().DeletePassword(ids, userInfo)
	if err != nil {
		return nil, err
	}

	return passwords, nil
}
