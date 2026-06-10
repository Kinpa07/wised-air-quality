package server

import (
	"net/http"

	loggerpkg "github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/router"
	"github.com/prometheus/client_golang/prometheus"
)

func NewRouter(serviceName string, prometheusRegisterer prometheus.Registerer, cfg *router.Config, loggerConfig *loggerpkg.Config, middlewares ...func(next http.Handler) http.Handler) *router.Router {
	return router.NewRouter(cfg, loggerConfig, prometheusRegisterer, serviceName, middlewares...)
}
