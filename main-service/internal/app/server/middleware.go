package server

import (
	"net/http"

	"github.com/SintroSecurity/go-libraries/db"
)

func CreateRootMiddleware(cfg *Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = db.NewContextWithDatabase(ctx, cfg.Database.WithContext(ctx))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
