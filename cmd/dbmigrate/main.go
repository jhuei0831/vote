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
    // 建立新的 FlagSet
    flags = flag.NewFlagSet("goose", flag.ExitOnError)
    // 設定遷移檔案的目錄
    dir = flags.String("dir", os.Getenv("GOOSE_MIGRATION_DIR"), "directory with migration files")
    // 決定migrate的動作是seed還是migration
    action = flags.String("action", "migrate", "action to perform: seed or migrate")
    // -no-versioning apply migration commands with no versioning, in file order, from directory pointed to
    noVersioning = flags.Bool("no-versioning", false, "apply migration commands with no versioning, in file order, from directory pointed to")
)

func main() {
    // 解析命令列參數
    err := flags.Parse(os.Args[2:])
    if err != nil {
        panic(err)
    }

    // 載入 .env 檔案
    if err := godotenv.Load(); err != nil {
        panic(err)
    }

    if *dir == "" {
        *dir = os.Getenv("GOOSE_MIGRATION_DIR")
    }
    
    // 動作，預設是 migrate，如果是seed，則次竟改成GOOSE_SEED_DIR
    if *action == "seed" {
        *dir = os.Getenv("GOOSE_SEED_DIR")
        *noVersioning = true
    } else if *action != "migrate" {
        slog.Error("Invalid action", "action", *action)
        return
    }

    // 取得解析後的參數
    args := flags.Args()
    slog.Info("Args are", "args", args, "dir", *dir, "action", *action)

    // 如果參數數量小於 1，顯示使用說明並返回
    if len(args) < 1 {
        flags.Usage()
        return
    }

    // 取得命令
    command := args[0]

    // 從環境變數取得資料庫設定
    dbConfig := os.Getenv("DB_CONFIG")
    // 初始化資料庫
    db, err := database.Initialize(dbConfig)
    if err != nil {
        panic(err)
    }

    // 取得 SQL 資料庫物件
    sqlDB, err := db.DB()
    if err != nil {
        panic(err)
    }

    // 確保在函數結束時關閉資料庫連線
    defer func() {
        if err := sqlDB.Close(); err != nil {
            panic(err)
        }
    }()

    // 準備命令的參數
    arguments := make([]string, 0)
    if len(args) > 1 {
        arguments = append(arguments, args[1:]...)
    }

    options := []goose.OptionsFunc{}
    if *noVersioning {
		options = append(options, goose.WithNoVersioning())
	}

    // 設定 goose 的資料庫方言
    if err := goose.SetDialect("mysql"); err != nil {
        panic(err)
    }

    // 執行 goose 命令
    if err := goose.RunWithOptionsContext(context.Background(), command, sqlDB, *dir, arguments, options...); err != nil {
        log.Fatalf("goose %v: %v", command, err)
    }
}