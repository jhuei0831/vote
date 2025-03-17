package seed

import (
	"context"
	"database/sql"

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
	passwordHash, err := SHA256Hasher.HashPassword("admin")
	if err != nil {
		transaction.Rollback()
		return err
	}
	// Create admin user
	user := model.User{
		Account:  "admin",
		Password: passwordHash,
		Email:    "admin@example.com",
	}
	err = transaction.Model(&model.User{}).Create(&user).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	// Create admin role
	enforcer := database.Enforcer
	_, err = enforcer.AddRoleForUser("admin", "ROLE_ADMIN")
	enforcer.AddPolicy("ROLE_ADMIN", "users", "read")
	enforcer.AddPolicy("ROLE_ADMIN", "users", "create")
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

	transaction.Migrator().DropTable("casbin_rule")

	return transaction.Commit().Error
}
