package migrations

import (
	"context"
	"database/sql"

	"vote/app/database"
	"vote/app/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable00001, downCreateUsersTable00001)
}

func upCreateUsersTable00001(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)

	if err != nil {
		return err
	}

	err = database.SqlSession.Migrator().CreateTable(&model.User{})

	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE OR REPLACE TRIGGER update_users_updated_at
		BEFORE UPDATE
		ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`)

	return err
}

func downCreateUsersTable00001(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return database.SqlSession.Migrator().DropTable(&model.User{})
}
