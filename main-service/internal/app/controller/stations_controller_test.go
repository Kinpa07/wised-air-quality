package controller

import (
	"context"
	"testing"
	"time"

	"go-service-skeleton/internal/app/display"
	"go-service-skeleton/internal/app/geo"
	"go-service-skeleton/internal/testutil"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"

	"gorm.io/gorm"
)

const (
	mitteClient  = "33333333-3333-3333-3333-333333333333"
	silentClient = "44444444-4444-4444-4444-444444444444"

	// Brandenburg Gate — lands in Mitte.
	mitteLat = 52.5163
	mitteLng = 13.3777
)

// stationsCtx augments the db context with a display config and the real district
// index, since GetStations reads both from context (the middleware isn't running
// in a controller test). expected = window/interval = 60/30 = 2, so seeding 1 or 2
// readings maps to ratios 0.5 / 1.0 without inserting hundreds of rows.
func stationsCtx(t *testing.T) (context.Context, *gorm.DB) {
	t.Helper()
	ctx, gdb := ctxWithDB(t)

	ctx = display.NewContext(ctx, &display.Config{
		ConnectionWindowMinutes: 60,
		ExpectedIntervalMinutes: 30,
		PoorConnectionThreshold: 0.8,
	})

	idx, err := geo.Load("../../../assets/berlin-districts.geojson")
	if err != nil {
		t.Fatalf("load districts: %v", err)
	}
	ctx = geo.NewContext(ctx, idx)

	return ctx, gdb
}

func findInData(t *testing.T, data []sensor_readings_collector_pkg.Station, id string) sensor_readings_collector_pkg.Station {
	t.Helper()
	for _, s := range data {
		if s.ID == id {
			return s
		}
	}
	t.Fatalf("station %s not in response", id)
	return sensor_readings_collector_pkg.Station{}
}

func Test_GetStations_EnrichesReportingStation(t *testing.T) {
	ctx, gdb := stationsCtx(t)
	testutil.SeedClientAt(t, gdb, mitteClient, mitteLat, mitteLng)

	now := time.Now().UTC()
	// Two readings inside the 60-min window => received 2, expected 2, ratio 1.0.
	testutil.SeedReading(t, gdb, mitteClient, 9, 18, now.Add(-20*time.Minute))
	testutil.SeedReading(t, gdb, mitteClient, 14, 22, now.Add(-5*time.Minute)) // newest

	res, errResp := GetStations(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}

	s := findInData(t, res.Data, mitteClient)

	if s.District != "Mitte" {
		t.Fatalf("district: want Mitte, got %q", s.District)
	}
	if s.PM25 == nil || *s.PM25 != 14 {
		t.Fatalf("pm2_5: want newest 14, got %v", s.PM25)
	}
	if s.Band == nil || *s.Band != sensor_readings_collector_pkg.AQIBandModerate {
		t.Fatalf("band: want Moderate, got %v", s.Band)
	}
	if s.Stability != 100 {
		t.Fatalf("stability: want 100, got %v", s.Stability)
	}
	if s.Connection != sensor_readings_collector_pkg.ConnectionGood {
		t.Fatalf("connection: want Good, got %v", s.Connection)
	}
}

func Test_GetStations_PoorConnection(t *testing.T) {
	ctx, gdb := stationsCtx(t)
	testutil.SeedClientAt(t, gdb, mitteClient, mitteLat, mitteLng)

	// One reading in the window => ratio 0.5, below the 0.8 threshold.
	testutil.SeedReading(t, gdb, mitteClient, 40, 60, time.Now().UTC().Add(-10*time.Minute))

	res, errResp := GetStations(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}

	s := findInData(t, res.Data, mitteClient)
	if s.Stability != 50 {
		t.Fatalf("stability: want 50, got %v", s.Stability)
	}
	if s.Connection != sensor_readings_collector_pkg.ConnectionPoor {
		t.Fatalf("connection: want Poor, got %v", s.Connection)
	}
}

func Test_GetStations_SilentStation(t *testing.T) {
	ctx, gdb := stationsCtx(t)
	testutil.SeedClientAt(t, gdb, silentClient, mitteLat, mitteLng) // enrolled, no readings

	res, errResp := GetStations(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}

	s := findInData(t, res.Data, silentClient)
	if s.PM25 != nil || s.PM10 != nil || s.MeasuredAt != nil || s.Band != nil {
		t.Fatalf("silent station: want nil reading fields, got %+v", s)
	}
	if s.Stability != 0 {
		t.Fatalf("stability: want 0, got %v", s.Stability)
	}
	if s.Connection != sensor_readings_collector_pkg.ConnectionPoor {
		t.Fatalf("connection: want Poor, got %v", s.Connection)
	}
	if s.District != "Mitte" {
		t.Fatalf("district: want Mitte even when silent, got %q", s.District)
	}
}

func Test_GetStations_EmptyFleet(t *testing.T) {
	ctx, _ := stationsCtx(t)

	res, errResp := GetStations(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	if len(res.Data) != 0 {
		t.Fatalf("want empty data, got %d rows", len(res.Data))
	}
}
