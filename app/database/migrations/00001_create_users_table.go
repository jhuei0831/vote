package migrations

import (
	"context"
	"database/sql"

	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable00001, downCreateUsersTable00001)
}

func upCreateUsersTable00001(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return database.SqlSession.Migrator().CreateTable(&model.User{})
}

func downCreateUsersTable00001(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.User{})
}
