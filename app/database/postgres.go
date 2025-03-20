package database

import (
	"time"
	"vote/app/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var (
	SqlSession *gorm.DB
	Adapter *gormadapter.Adapter
	Enforcer *casbin.Enforcer
)

func Initialize(dbConfig string) (*gorm.DB, error) {
	newLogger := logger.New(
		utils.Logger(),
		logger.Config{
			SlowThreshold:              time.Second,   // Slow SQL threshold
			LogLevel:                   logger.Info, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,           // Don't include params in the SQL log
			Colorful:                  false,          // Disable color
		},
	)
	var err error
	SqlSession, err = gorm.Open(postgres.Open(dbConfig), &gorm.Config{
		Logger: newLogger,
	})

	return SqlSession, err
}

func Rbac() (*gormadapter.Adapter, *casbin.Enforcer, error) {
	var err error

	Adapter, err = gormadapter.NewAdapterByDB(SqlSession)
	if err != nil {
		return nil, nil, err
	}
	
	Enforcer, err = casbin.NewEnforcer("app/config/rbac.conf", Adapter)
	if err != nil {
		return nil, nil, err
	}

	Enforcer.EnableLog(true)
	
	return Adapter, Enforcer, err
}