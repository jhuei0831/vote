package seed

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/pressly/goose/v3"

	"vote/app/database"
	"vote/app/model"
	"vote/app/utils"
)

func init() {
	goose.AddMigrationContext(upSeedUsersTable, downSeedUsersTable)
}

func upSeedUsersTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	transaction := database.SqlSession.Begin()

	var SHA256Hasher utils.SHA256Hasher
	passwordHash, err := SHA256Hasher.HashPassword("password")
	if err != nil {
		transaction.Rollback()
		return err
	}
	// Create admin user
	users := []model.User{
		{
			Account:  "admin",
			Password: passwordHash,
			Email:    "admin@example.com",
		},
		{
			Account:  "creator",
			Password: passwordHash,
			Email:    "creator@example.com",
		},
	}
	err = transaction.Model(&model.User{}).Create(&users).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Create admin role
	_, enforcer, err := database.Rbac()
	if err != nil {
		transaction.Rollback()
	}
	for _, user := range users {
		role := strings.ToUpper(user.Account)
		userId := strconv.FormatUint(user.ID, 10)
		_, err = enforcer.AddRoleForUser(userId, role)
		enforcer.AddPolicy(role, "user", "create")
		enforcer.AddPolicy(role, "user", "read")
		enforcer.AddPolicy(role, "user", "update")
		enforcer.AddPolicy(role, "user", "delete")
		enforcer.AddPolicy(role, "vote", "create")
		enforcer.AddPolicy(role, "vote", "read")
		enforcer.AddPolicy(role, "vote", "update")
		enforcer.AddPolicy(role, "vote", "delete")
	}
	
	if err != nil {
		transaction.Rollback()
	}

	return transaction.Commit().Error
}

func downSeedUsersTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	transaction := database.SqlSession.Begin()

	err := transaction.Model(&model.User{}).Unscoped().Where("account = ?", "admin").Delete(&model.User{}).Error

	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit().Error
}
