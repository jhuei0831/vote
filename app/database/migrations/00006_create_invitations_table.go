package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateInvitationsTable00006, downCreateInvitationsTable00006)
}

func upCreateInvitationsTable00006(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return database.SqlSession.Migrator().CreateTable(&model.Invitation{})
}

func downCreateInvitationsTable00006(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Invitation{})
}
