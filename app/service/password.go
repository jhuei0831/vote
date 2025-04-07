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

// SelectPassword 根據提供的投票ID，檢索所有密碼。
func (p PasswordService) SelectPassword(voteId uuid.UUID, passwordQuery model.PasswordQuery) ([]model.Password, int64, error) {
	var passwords []model.Password
	var total int64
	query := database.SqlSession.Model(&passwords).Where("vote_id = ?", voteId)

	if passwordQuery.Status {
		query = query.Where("status = ?", passwordQuery.Status)
	}

	// 設定查詢條件
	page := passwordQuery.Page
	size := passwordQuery.Size

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

	err = query.Find(&passwords).Error
	if err != nil {
		return nil, 0, err
	}

	return passwords, total, nil
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