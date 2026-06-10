package app

import (
	"github.com/SintroSecurity/go-libraries/auth"
	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/router"
	"github.com/SintroSecurity/go-libraries/worker"
)

type Config struct {
	Server *router.Config `mapstructure:"SERVER"`
	Auth   *auth.Config   `mapstructure:"AUTH"`
	Logger *logger.Config `mapstructure:"LOGGER"`
	DB     *db.Config     `mapstructure:"DB"`
	Worker *worker.Config `mapstructure:"WORKER"`
}
