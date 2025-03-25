package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateQuestionsTable00004, downCreateQuestionsTable00004)
}

func upCreateQuestionsTable00004(ctx context.Context, tx *sql.Tx) error {
	return database.SqlSession.Migrator().CreateTable(&model.Question{})
}

func downCreateQuestionsTable00004(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Question{})
}
