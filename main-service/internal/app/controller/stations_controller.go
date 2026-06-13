package controller

import (
	"context"
	"go-service-skeleton/internal/app/database"
	"go-service-skeleton/internal/app/display"
	"go-service-skeleton/internal/app/geo"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
	"time"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
)

// latestPerStationQuery left-joins every client to its cached latest reading
// (maintained on ingest), avoiding a scan of the whole readings history. The
// deleted_at filter is manual: raw SQL bypasses GORM's soft-delete scope.
const latestPerStationQuery = `
SELECT c.id, c.type, c.latitude, c.longitude,
       lr.pm2_5, lr.pm10, lr.measured_at
FROM clients c
LEFT JOIN latest_readings lr ON lr.client_id = c.id
WHERE c.deleted_at IS NULL
`

// receivedCountQuery counts each station's readings since the cutoff. Stations
// with none are absent from the result, so the Go-side merge reads them as zero.
const receivedCountQuery = `
SELECT client_id, COUNT(*) AS received
FROM readings
WHERE measured_at >= ?
GROUP BY client_id
`

type receivedRow struct {
	ClientID string
	Received int
}

func GetStations(ctx context.Context) (*sensor_readings_collector_pkg.GetStationsResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)
	districts := geo.FromContext(ctx)

	var stations []database.StationRow

	result := gdb.Raw(latestPerStationQuery).Scan(&stations)

	if result.Error != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	cfg := display.FromContext(ctx)
	cutoff := time.Now().UTC().Add(-time.Duration(cfg.ConnectionWindowMinutes) * time.Minute)

	var received []receivedRow

	result = gdb.Raw(receivedCountQuery, cutoff).Scan(&received)
	if result.Error != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	countByClient := make(map[string]int)
	for _, c := range received {
		countByClient[c.ClientID] = c.Received
	}

	data := make([]sensor_readings_collector_pkg.Station, len(stations))
	expected := 0
	if cfg.ExpectedIntervalMinutes > 0 {
		expected = cfg.ConnectionWindowMinutes / cfg.ExpectedIntervalMinutes

	}
	for i, s := range stations {
		received := countByClient[s.ID]

		ratio := 0.0
		if expected > 0 {
			ratio = float64(received) / float64(expected)
		}

		band := sensor_readings_collector_pkg.BandFor(s.PM25)
		district := districts.DistrictForClient(s.ID, s.Latitude, s.Longitude)

		data[i] = sensor_readings_collector_pkg.Station{
			ID:         s.ID,
			Latitude:   s.Latitude,
			Longitude:  s.Longitude,
			PM25:       s.PM25,
			PM10:       s.PM10,
			MeasuredAt: s.MeasuredAt,
			Band:       band,
			District:   district,
			Stability:  ratio * 100,
			Connection: sensor_readings_collector_pkg.QualityFor(ratio, cfg.PoorConnectionThreshold),
		}
	}

	return &sensor_readings_collector_pkg.GetStationsResponse{Data: data}, nil
}
