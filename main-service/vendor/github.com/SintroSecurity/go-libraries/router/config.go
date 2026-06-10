package router

import (
	"github.com/SintroSecurity/go-libraries/router/ratelimiter"
	"github.com/go-chi/cors"
)

// Config holds the configuration for the router

type Config struct {
	ListenURL   string                   `mapstructure:"LISTENURL"`
	CORS        *CORSConfig              `mapstructure:"CORS"`
	RateLimiter *ratelimiter.RateLimiter `mapstructure:"RATELIMITER"`
}
type CORSConfig struct {
	// CORSAllowedOrigins is a list of allowed origins for CORS requests
	// If empty and EnableCORS is true, defaults to ["*"]
	AllowedOrigins []string `mapstructure:"ALLOWEDORIGINS"`
	// CORSAllowedMethods is a list of allowed HTTP methods for CORS requests
	// If empty and EnableCORS is true, defaults to ["GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"]
	AllowedMethods []string `mapstructure:"ALLOWEDMETHODS"`
	// CORSAllowedHeaders is a list of allowed headers for CORS requests
	// If empty and EnableCORS is true, defaults to ["Accept", "Authorization", "Content-Type", "X-Request-ID"]
	AllowedHeaders []string `mapstructure:"ALLOWEDHEADERS"`
	// CORSExposedHeaders is a list of headers exposed to the client
	// If empty and EnableCORS is true, defaults to ["Link", "X-Request-ID"]
	ExposedHeaders []string `mapstructure:"EXPOSEDHEADERS"`
	// CORSAllowCredentials indicates whether credentials are allowed
	// Defaults to false
	AllowCredentials bool `mapstructure:"ALLOWCREDENTIALS"`
	// CORSMaxAge indicates how long preflight requests can be cached (in seconds)
	// Defaults to 300
	MaxAge int `mapstructure:"MAXAGE"`
}

// newOptionsWithCORS creates a default CORS configuration with CORS enabled and sensible defaults
func newOptionsWithCORS() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}
}

func NewCORSOptions(c *CORSConfig) cors.Options {
	defaultCORSOptions := newOptionsWithCORS()

	if c.AllowedOrigins != nil {
		defaultCORSOptions.AllowedOrigins = c.AllowedOrigins
	}
	if c.AllowedMethods != nil {
		defaultCORSOptions.AllowedMethods = c.AllowedMethods
	}
	if c.AllowedHeaders != nil {
		defaultCORSOptions.AllowedHeaders = c.AllowedHeaders
	}
	if c.ExposedHeaders != nil {
		defaultCORSOptions.ExposedHeaders = c.ExposedHeaders
	}
	defaultCORSOptions.AllowCredentials = c.AllowCredentials
	if c.MaxAge != 0 {
		defaultCORSOptions.MaxAge = c.MaxAge
	}

	return defaultCORSOptions

}
