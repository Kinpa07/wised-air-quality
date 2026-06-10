package db

import (
	"database/sql"
	"fmt"
)

func GetDSN(cfg *Config) string {
	switch cfg.Dialect {
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Port, cfg.Postgres.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)
	case "sqlite3":
		return fmt.Sprintf("file:%s?_fk=yes", cfg.SQLite.Filename)
	case "clickhouse":
		return fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s", cfg.Clickhouse.Username, cfg.Clickhouse.Password, cfg.Clickhouse.Host, cfg.Clickhouse.Port, cfg.Clickhouse.Database)
	default:
		panic("invalid dialect")
	}
}

// InitOpenedDB initializes early opened *sql.DB instance.
func InitOpenedDB(db *sql.DB, cfg *Config) error {
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

	return db.Ping()
}
