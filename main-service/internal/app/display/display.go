// Package display holds the operator-tunable fleet thresholds and the context
// plumbing to carry them per request. It is a leaf package so both app (loads it)
// and controller (reads it) can import it without an import cycle.
package display

import "context"

type Config struct {
	ActiveWindowMinutes     int     `mapstructure:"ACTIVE_WINDOW_MINUTES"`     // reading newer than this => station "active"
	ConnectionWindowMinutes int     `mapstructure:"CONNECTION_WINDOW_MINUTES"` // lookback for the received/expected ratio
	ExpectedIntervalMinutes int     `mapstructure:"EXPECTED_INTERVAL_MINUTES"` // assumed cadence; expected = window / interval
	PoorConnectionThreshold float64 `mapstructure:"POOR_CONNECTION_THRESHOLD"` // ratio below this => "poor connection"
	DefaultPageSize         int     `mapstructure:"DEFAULT_PAGE_SIZE"`         // cursor page size when ?limit= is absent
}

type ctxKey struct{}

func NewContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, ctxKey{}, cfg)
}

func FromContext(ctx context.Context) *Config {
	cfg, _ := ctx.Value(ctxKey{}).(*Config)
	return cfg
}
