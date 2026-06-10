package ratelimiter

import (
	"golang.org/x/time/rate"
	"time"
)

type Config struct {
	Limit rate.Limit    `mapstructure:"LIMIT"`
	Burst int           `mapstructure:"BURST"`
	TTL   time.Duration `mapstructure:"TTL"`
}
