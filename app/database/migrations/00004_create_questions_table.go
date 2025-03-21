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
	err := database.SqlSession.Migrator().CreateTable(&model.Question{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE OR REPLACE TRIGGER update_questions_updated_at
		BEFORE UPDATE
		ON questions
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`)

	return err
}

func downCreateQuestionsTable00004(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.Question{})
}
