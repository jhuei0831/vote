package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateBallotsTable00007, downCreateBallotsTable00007)
}

func upCreateBallotsTable00007(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return database.SqlSession.Migrator().CreateTable(&model.Ballot{})
}

func downCreateBallotsTable00007(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Ballot{})
}
