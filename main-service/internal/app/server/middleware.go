package server

import (
	"net/http"

	"go-service-skeleton/internal/app/display"
	"go-service-skeleton/internal/app/geo"

	"github.com/SintroSecurity/go-libraries/db"
)

func CreateRootMiddleware(cfg *Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = db.NewContextWithDatabase(ctx, cfg.Database.WithContext(ctx))
			ctx = display.NewContext(ctx, cfg.Display)
			ctx = geo.NewContext(ctx, cfg.Districts)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
