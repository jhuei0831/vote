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
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	flags := flag.NewFlagSet("goose", flag.ExitOnError)
    dir   := flags.String("dir", os.Getenv("GOOSE_MIGRATION_DIR"), "directory with migration files")
	
    err := flags.Parse(os.Args[2:])
    if err != nil {
        panic(err)
    }

    args := flags.Args()
    slog.Info("Args are", "args", args, "dir", *dir)

    if len(args) < 1 {
        flags.Usage()
        return
    }

    command := args[0]

	dbConfig := os.Getenv("DB_CONFIG")
	db, err := database.Initialize(dbConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

    defer func() {
        if err := sqlDB.Close(); err != nil {
            panic(err)
        }
    }()

    arguments := make([]string, 0)
    if len(args) > 1 {
        arguments = append(arguments, args[3:]...)
    }

	if err := goose.SetDialect("mysql"); err != nil {
        panic(err)
    }

    if err := goose.RunContext(context.Background(), command, sqlDB, *dir, arguments...); err != nil {
        log.Fatalf("goose %v: %v", command, err)
    }
}