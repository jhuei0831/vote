package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreatePollOptionsTable00005, downCreatePollOptionsTable00005)
}

func upCreatePollOptionsTable00005(ctx context.Context, tx *sql.Tx) error {
	return database.SqlSession.Migrator().CreateTable(&model.PollOption{})
}

func downCreatePollOptionsTable00005(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.PollOption{})
}
