package database

import (
	// "log"
	// "os"
	"time"
	"vote/app/middleware"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


var (
	SqlSession *gorm.DB
)

func Initialize(dbConfig string) (*gorm.DB, error) {
	newLogger := logger.New(
		// log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		middleware.Logger(),
		logger.Config{
			SlowThreshold:              time.Second,   // Slow SQL threshold
			LogLevel:                   logger.Info, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  false,          // Disable color
		},
	)
	var err error
	SqlSession, err = gorm.Open(mysql.Open(dbConfig), &gorm.Config{
		Logger: newLogger,
	})

	return SqlSession, err
}