package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"

	"vote/app/database"
	"vote/app/model"
)

func init() {
	goose.AddMigrationContext(upCreateVotesTable00003, downCreateVotesTable00003)
}

func upCreateVotesTable00003(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return database.SqlSession.Migrator().CreateTable(&model.Vote{})
}

func downCreateVotesTable00003(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Vote{})
}
