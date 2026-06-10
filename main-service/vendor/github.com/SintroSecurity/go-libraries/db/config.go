package db

import "time"

type Config struct {
	Dialect               string        `mapstructure:"DIALECT"`
	MaxOpenConnections    int           `mapstructure:"MAXOPENCONNECTIONS"`
	MaxIdleConnections    int           `mapstructure:"MAXIDLECONNECTIONS"`
	ConnectionMaxLifetime time.Duration `mapstructure:"CONNECTIONMAXLIFETIME"`
	Postgres              *Postgres     `mapstructure:"POSTGRES"`
	MySQL                 *MySQL        `mapstructure:"MYSQL"`
	SQLite                *SQLite       `mapstructure:"SQLITE"`
	Clickhouse            *Clickhouse   `mapstructure:"CLICKHOUSE"`
}

type Postgres struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Database string `mapstructure:"DATABASE"`
	SSLMode  string `mapstructure:"SSLMODE"`
}

type MySQL struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Database string `mapstructure:"DATABASE"`
}

type SQLite struct {
	Filename string `mapstructure:"FILENAME"`
}

type Clickhouse struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	Database string `mapstructure:"DATABASE"`
}
