package service

import (
	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

	"github.com/google/uuid"
)

type PasswordService struct {
}

func NewPasswordService() PasswordService {
	return PasswordService{}
}

// SelectOnePassword 根據提供的投票ID和密碼，檢查密碼是否存在。
func (p PasswordService) SelectOnePassword(voteId uuid.UUID, password string) (*model.Password, error) {
	passwordModel := model.Password{}
	err := database.SqlSession.
		Where("vote_id = ? AND password = ?", voteId, password).
		First(&passwordModel).
		Error
		
	if err != nil {
		return nil, err
	}
	
	return &passwordModel, nil
}

// CreatePassword 建立可以加解密的密碼
func (p PasswordService) CreatePassword(voteId uuid.UUID, number int, length int, format string) error {
	passwordUtil := &utils.Password{}
	// 生成密碼
	passwords, err := passwordUtil.GeneratePassword(number, length, format)
	if err != nil {
		return err
	}

	// 將密碼加密
	passwordModels := make([]model.Password, len(passwords))
	for i, password := range passwords {
		passwordEncrypt, err := passwordUtil.Encrypt(password)
		if err != nil {
			return err
		}
		passwordModels[i] = model.Password{
			VoteID:   voteId,
			Password: passwordEncrypt,
		}
	}

	// 使用transaction，將密碼存入資料庫
	transaction := database.SqlSession.Begin()
	err = transaction.CreateInBatches(&passwordModels, 100).Error

	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit().Error
}