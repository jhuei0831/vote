package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"vote/app/database"
	_ "vote/app/database/migrations"
	_ "vote/app/database/seed"
)

var (
	// Create a new FlagSet
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	// Set the directory for migration files
	dir = flags.String("dir", os.Getenv("GOOSE_MIGRATION_DIR"), "directory with migration files")
	// Determine whether the migrate action is seed or migration
	action = flags.String("action", "migrate", "action to perform: seed or migrate")
	// -no-versioning apply migration commands with no versioning, in file order, from directory pointed to
	noVersioning = flags.Bool("no-versioning", false, "apply migration commands with no versioning, in file order, from directory pointed to")
)

func main() {
	// Parse command line arguments
	err := flags.Parse(os.Args[2:])
	if err != nil {
		panic(err)
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if *dir == "" {
		*dir = os.Getenv("GOOSE_MIGRATION_DIR")
	}

	// Action defaults to migrate; if seed, change directory to GOOSE_SEED_DIR
	if *action == "seed" {
		*dir = os.Getenv("GOOSE_SEED_DIR")
		*noVersioning = true
	} else if *action != "migrate" {
		slog.Error("Invalid action", "action", *action)
		return
	}

	// Get parsed arguments
	args := flags.Args()
	slog.Info("Args are", "args", args, "dir", *dir, "action", *action)

	// If number of arguments is less than 1, show usage and return
	if len(args) < 1 {
		flags.Usage()
		return
	}

	// Get the command
	command := args[0]

	// Get database configuration from environment variables
	dbConfig := database.DbConfig()
	// Initialize database
	db, err := database.Initialize(dbConfig)
	if err != nil {
		panic(err)
	}

	// Get SQL database object
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// Ensure database connection is closed when function ends
	defer func() {
		if err := sqlDB.Close(); err != nil {
			panic(err)
		}
	}()

	// Prepare command arguments
	arguments := make([]string, 0)
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	options := []goose.OptionsFunc{}
	if *noVersioning {
		options = append(options, goose.WithNoVersioning())
	}

	// Set goose database dialect
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	// Execute goose command
	if err := goose.RunWithOptionsContext(context.Background(), command, sqlDB, *dir, arguments, options...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
