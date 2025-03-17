package service

import (
	"fmt"

	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

	"github.com/sirupsen/logrus"
)

var UserFields = []string{"id", "account", "email"}

func SelectOneUsers(id int64) (*model.User, error) {
	userOne := &model.User{}
	err := database.SqlSession.Select(UserFields).Where("id=?", id).First(&userOne).Error
	if err != nil {
		return nil, err
	} else {
		return userOne, nil
	}
}

func RegisterOneUser(account string, password string, email string) error {
	if !CheckOneUser(account) {
		return fmt.Errorf("user exists")
	}
	var SHA256Hasher utils.SHA256Hasher
	passwordHash, err := SHA256Hasher.HashPassword(password)
	if err != nil {
		return err
	}

	user := model.User{
		Account:  account,
		Password: passwordHash,
		Email:    email,
	}

	insertErr := database.SqlSession.Model(&model.User{}).Create(&user).Error
	utils.Logger().WithFields(logrus.Fields{
		"name": "RegisterOneUser",
	}).Error("error: ", insertErr)
	return insertErr
}

// CheckOneUser 根據提供的帳號檢查用戶是否存在。
// 如果用戶存在，返回 true；否則返回 false。
// 參數:
//   - account: 用戶的帳號。
//
// 返回值:
//   - bool: 如果用戶存在返回 true，否則返回 false。
func CheckOneUser(account string) bool {
	result := false
	var user model.User

	dbResult := database.SqlSession.Where("account = ?", account).Find(&user)
	if dbResult.Error != nil {
		fmt.Printf("Get User Info Failed:%v\n", dbResult.Error)
	} else {
		result = true
	}
	return result
}

// LoginOneUser 根據提供的帳號和密碼登錄用戶。
// 它首先使用 SHA256Hasher 對密碼進行哈希處理，然後比較哈希值以驗證密碼。
// 如果密碼驗證成功，則從資料庫中查找對應帳號的用戶資料。
// 參數:
//   - account: 用戶帳號
//   - password: 用戶密碼
//
// 返回值:
//   - *model.User: 如果登錄成功，返回用戶資料
//   - error: 如果登錄失敗，返回錯誤信息
func LoginOneUser(account string, password string) (*model.User, error) {
	var user model.User
	// Hash password
	var SHA256Hasher utils.SHA256Hasher
	passwordHash, err := SHA256Hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}
	// Check password
	if !SHA256Hasher.ComparePassword(password, passwordHash) {
		return nil, fmt.Errorf("password error")
	}

	dbResult := database.SqlSession.Where("account = ?", account).First(&user)

	if dbResult.Error != nil {
		return nil, dbResult.Error
	} else {
		return &user, nil
	}
}