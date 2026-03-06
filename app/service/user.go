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

func (u UserService) GetUserById(id uint64) (*model.User, error) {
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

// CheckAccountExist checks if a user exists based on the provided account.
// Returns true if the user exists, false otherwise.
// Parameters:
//   - account: The user's account name.
//
// Returns:
//   - bool: Returns true if the user exists, false otherwise.
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

// LoginOneUser authenticates a user based on the provided account and password.
// It retrieves the user from the database, then compares the hashed password using SHA256Hasher.
// If the password verification succeeds, it returns the user data.
// Parameters:
//   - account: The user's account name.
//   - password: The user's password.
//
// Returns:
//   - *model.User: Returns the user data if login is successful.
//   - error: Returns an error message if login fails.
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
