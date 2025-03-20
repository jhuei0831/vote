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
	err := database.SqlSession.Migrator().CreateTable(&model.Vote{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE OR REPLACE TRIGGER update_votes_updated_at
		BEFORE UPDATE
		ON votes
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`)

	return err
}

func downCreateVotesTable00003(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Vote{})
}
