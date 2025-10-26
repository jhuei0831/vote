package service

import (
	"fmt"

	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"

	"github.com/sirupsen/logrus"
)

type UserService struct {
}

func NewUserService() UserService {
	return UserService{}
}

func (u UserService) GetUserById(id int64) (*model.User, error) {
	user := &model.User{}
	err := database.SqlSession.Select([]string{"id", "account", "email"}).Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

func (u UserService) GetUsers() ([]*model.User, error) {
	var users []*model.User

	err := database.SqlSession.Select([]string{"id", "account", "email"}).Find(&users).Error

	if err != nil {
		return nil, err
	} else {
		return users, err
	}
}

func (u UserService) CreateUser(input model.UserCreate) (*model.User, error) {
	if u.CheckAccountExist(input.Account) {
		return nil, fmt.Errorf("account exists")
	}
	if u.CheckEmailExist(input.Email) {
		return nil, fmt.Errorf("email exists")
	}
	var SHA256Hasher utils.SHA256Hasher
	passwordHash, err := SHA256Hasher.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Account:  input.Account,
		Password: passwordHash,
		Email:    input.Email,
	}

	insertErr := database.SqlSession.Model(&model.User{}).Create(&user).Error
	utils.Logger().WithFields(logrus.Fields{
		"name": "CreateUser",
	}).Error("error: ", insertErr)

	return &user, insertErr
}

// CheckAccountExist 根據提供的帳號檢查用戶是否存在。
// 如果用戶存在，返回 true；否則返回 false。
// 參數:
//   - account: 用戶的帳號。
//
// 返回值:
//   - bool: 如果用戶存在返回 true，否則返回 false。
func (u UserService) CheckAccountExist(account string) bool {
	var user model.User
	dbResult := database.SqlSession.Model(&model.User{}).Select("id").Where("account = ?", account).Limit(1).Find(&user)

	if dbResult.Error != nil {
		utils.Logger().WithFields(logrus.Fields{
			"name":    "CheckAccountExist",
			"account": account,
		}).Error(dbResult.Error)
		return false
	}

	return dbResult.RowsAffected > 0
}

func (u UserService) CheckEmailExist(email string) bool {
	var user model.User
	dbResult := database.SqlSession.Model(&model.User{}).Select("id").Where("email = ?", email).Limit(1).Find(&user)

	if dbResult.Error != nil {
		utils.Logger().WithFields(logrus.Fields{
			"name":  "CheckEmailExist",
			"email": email,
		}).Error(dbResult.Error)
		return false
	}

	return dbResult.RowsAffected > 0
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
func (u UserService) LoginOneUser(account string, password string) (*model.User, error) {
	var user model.User
	var SHA256Hasher utils.SHA256Hasher

	// Get user info
	dbResult := database.SqlSession.Where("account = ?", account).First(&user)

	// Check password
	if !SHA256Hasher.ComparePassword(password, user.Password) {
		return nil, fmt.Errorf("password error")
	}

	if dbResult.Error != nil {
		return nil, dbResult.Error
	} else {
		return &user, nil
	}
}
