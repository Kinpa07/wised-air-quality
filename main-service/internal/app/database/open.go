package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SintroSecurity/go-libraries/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Open(ctx context.Context, cfg *db.Config) (*gorm.DB, error) {
	var gormDB *gorm.DB
	var err error

	dsn := db.GetDSN(cfg)

	var dialector gorm.Dialector

	switch cfg.Dialect {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite3":
		dialector = sqlite.Open(dsn)
	}

	jsonLogger := New(ctx, log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		Config{
			SlowThreshold:             time.Second,       // Slow SQL threshold
			LogLevel:                  gormLogger.Silent, // Log level
			IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,             // Disable color
		})

	if gormDB, err = gorm.Open(dialector, &gorm.Config{
		Logger:         jsonLogger,
		TranslateError: true,
	}); err != nil {
		return nil, fmt.Errorf("open database error: %w", err)
	}

	gormDB = gormDB.Session(&gorm.Session{Context: ctx})

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("can't get sql.DB with error: %w", err)
	}

	if err := db.InitOpenedDB(sqlDB, cfg); err != nil {
		return nil, fmt.Errorf("can't initialize DB with error: %w", err)
	}

	return gormDB, nil
}
