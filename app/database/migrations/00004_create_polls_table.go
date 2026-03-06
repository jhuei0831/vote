package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreatePollsTable00004, downCreatePollsTable00004)
}

func upCreatePollsTable00004(ctx context.Context, tx *sql.Tx) error {
	return database.SqlSession.Migrator().CreateTable(&model.Poll{})
}

func downCreatePollsTable00004(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Poll{})
}
