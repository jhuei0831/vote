package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateInvitationUsagesTable00009, downCreateInvitationUsagesTable00009)
}

func upCreateInvitationUsagesTable00009(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return database.SqlSession.Migrator().CreateTable(&model.InvitationUsage{})
}

func downCreateInvitationUsagesTable00009(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.InvitationUsage{})
}
