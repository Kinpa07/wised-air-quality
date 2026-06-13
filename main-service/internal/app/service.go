package app

import (
	"context"
	"go-service-skeleton/internal/app/database"
	"go-service-skeleton/internal/app/geo"
	server2 "go-service-skeleton/internal/app/server"
	"go-service-skeleton/internal/app/worker"
	"net/http"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/metrics"
	"github.com/SintroSecurity/go-libraries/router"
	"github.com/prometheus/client_golang/prometheus"
)

const ServiceName = "sensor-readings-collector"

const districtsPath = "assets/berlin-districts.geojson"

func InitRouterAndRoutingTable(ctx context.Context, prometheusRegisterer prometheus.Registerer, cfg *Config, rootMiddleware func(next http.Handler) http.Handler) (*router.Router, error) {
	r := server2.NewRouter(ServiceName, prometheusRegisterer, cfg.Server, cfg.Logger, rootMiddleware)
	server2.InitRoutingTable(ctx, r)
	return r, nil
}

func InitMetrics() prometheus.Registerer {
	metrics := metrics.NewMetrics()
	return metrics
}

func InitWorker(ctx context.Context, cfg *Config) error {
	return worker.CreateAndRegisterHandlers(ctx, cfg.Worker, ServiceName)
}

func StartService(ctx context.Context, cfg *Config) error {
	//It's a good idea to init logger on the top of this func, because once set in context will be passed down to each request
	l := logger.New(ctx, cfg.Logger)
	ctx = logger.NewContextWithLogger(ctx, l)

	gdb, err := database.Open(ctx, cfg.DB)
	if err != nil {
		return err
	}
	if err := database.MigrateDatabase(gdb); err != nil {
		return err
	}
	ctx = db.NewContextWithDatabase(ctx, gdb)

	districts, err := geo.Load(districtsPath)
	if err != nil {
		return err
	}

	registry := InitMetrics()
	ctx = metrics.NewContextWithPrometheusRegistry(ctx, registry)

	err = InitWorker(ctx, cfg)
	if err != nil {
		return err
	}

	rootMiddleware := server2.CreateRootMiddleware(&server2.Config{
		Database:  gdb,
		Display:   cfg.Display,
		Districts: districts,
	})

	r, err := InitRouterAndRoutingTable(ctx, registry, cfg, rootMiddleware)
	if err != nil {
		return err
	}
	r.Start()
	return nil
}
