package migrations

import (
	"context"
	"database/sql"
	"os"

	"github.com/pressly/goose/v3"
	"vote/app/database"
	"vote/app/model"
)

func init() {
	goose.AddMigrationContext(upCreateUserTable, downCreateUserTable)
}

func upCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	db, err := database.Initialize(os.Getenv("DB_CONFIG"))
	if err != nil {
		return err
	}
	// This code is executed when the migration is applied.
	return db.Migrator().CreateTable(&model.User{})
}

func downCreateUserTable(ctx context.Context, tx *sql.Tx) error {
	db, err := database.Initialize(os.Getenv("DB_CONFIG"))
	if err != nil {
		return err
	}
	// This code is executed when the migration is rolled back.
	return db.Migrator().DropTable(&model.User{})
}
