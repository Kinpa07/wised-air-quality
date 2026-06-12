package app

import (
	"github.com/SintroSecurity/go-libraries/auth"
	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/router"
	"github.com/SintroSecurity/go-libraries/worker"
)

type Config struct {
	Server  *router.Config `mapstructure:"SERVER"`
	Auth    *auth.Config   `mapstructure:"AUTH"`
	Logger  *logger.Config `mapstructure:"LOGGER"`
	DB      *db.Config     `mapstructure:"DB"`
	Worker  *worker.Config `mapstructure:"WORKER"`
	Display *DisplayConfig `mapstructure:"DISPLAY"`
}

// DisplayConfig holds the operator-tunable fleet thresholds the display API reads at request time.
type DisplayConfig struct {
	ActiveWindowMinutes     int     `mapstructure:"ACTIVE_WINDOW_MINUTES"`     // reading newer than this => station "active"
	ConnectionWindowMinutes int     `mapstructure:"CONNECTION_WINDOW_MINUTES"` // lookback for the received/expected ratio
	ExpectedIntervalMinutes int     `mapstructure:"EXPECTED_INTERVAL_MINUTES"` // assumed cadence; expected = window / interval
	PoorConnectionThreshold float64 `mapstructure:"POOR_CONNECTION_THRESHOLD"` // ratio below this => "poor connection"
	DefaultPageSize         int     `mapstructure:"DEFAULT_PAGE_SIZE"`         // cursor page size when ?limit= is absent
}
