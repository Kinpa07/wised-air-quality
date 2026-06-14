package controller

import (
	"context"
	"go-service-skeleton/internal/app/database"
	"go-service-skeleton/internal/app/display"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
	"time"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
)

func GetStats(ctx context.Context) (*sensor_readings_collector_pkg.GetStatsResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)
	cfg := display.FromContext(ctx)

	var stations []database.StationRow

	result := gdb.Raw(latestPerStationQuery).Scan(&stations)
	if result.Error != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	activeCutoff := time.Now().UTC().Add(-time.Duration(cfg.ActiveWindowMinutes) * time.Minute)

	countByClient, expected, err := connectionRatios(gdb, cfg)
	if err != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	ratioSum := 0.0
	activeCount := 0
	pmSum := 0.0
	poorCount := 0
	for _, s := range stations {
		received := countByClient[s.ID]
		ratio := 0.0
		if expected > 0 {
			ratio = float64(received) / float64(expected)
		}
		ratioSum += ratio
		// The LEFT JOIN fills all cache columns together, so a non-nil MeasuredAt
		// guarantees PM25 is non-nil — the deref below is safe.
		if s.MeasuredAt != nil && s.MeasuredAt.After(activeCutoff) {
			activeCount++
			pmSum += *s.PM25
		}

		if ratio < cfg.PoorConnectionThreshold {
			poorCount++
		}
	}

	var avg *float64
	if activeCount > 0 {
		v := pmSum / float64(activeCount)
		avg = &v
	}

	stability := 0.0
	if len(stations) > 0 {
		stability = ratioSum / float64(len(stations)) * 100
	}

	return &sensor_readings_collector_pkg.GetStatsResponse{
		Data: sensor_readings_collector_pkg.Stats{
			ActiveSensors:    activeCount,
			AvgPM25:          avg,
			PoorConnection:   poorCount,
			NetworkStability: stability,
		},
	}, nil
}
