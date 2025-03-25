package migrations

import (
	"context"
	"database/sql"
	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCandidatesTable00005, downCreateCandidatesTable00005)
}

func upCreateCandidatesTable00005(ctx context.Context, tx *sql.Tx) error {
	return database.SqlSession.Migrator().CreateTable(&model.Candidate{})
}

func downCreateCandidatesTable00005(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Candidate{})
}
