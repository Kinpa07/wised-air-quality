// Package testutil holds shared test fixtures for the in-memory SQLite harness
// reused across the controller, server, and display data-layer test suites.
package testutil

import (
	"fmt"
	"testing"
	"time"

	"go-service-skeleton/internal/app/database"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// OpenMemDB opens a migrated in-memory SQLite pinned to a single connection.
// It returns an error instead of failing a test, so callers without a
// *testing.T (notably TestMain, which only has *testing.M) can use it too.
func OpenMemDB() (*gorm.DB, error) {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("open in-memory sqlite: %w", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1) // one connection keeps the :memory: db alive for the run

	if err := database.Migrate(gdb); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return gdb, nil
}

// NewMemDB is the *testing.T/*testing.B convenience wrapper over OpenMemDB:
// it fails the test on error and closes the connection via t.Cleanup.
func NewMemDB(tb testing.TB) *gorm.DB {
	tb.Helper()

	gdb, err := OpenMemDB()
	if err != nil {
		tb.Fatalf("new mem db: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		tb.Fatalf("get sql.DB: %v", err)
	}
	tb.Cleanup(func() { _ = sqlDB.Close() })

	return gdb
}

// SeedClient inserts an enrolled station at a fixed Berlin coordinate (Mitte).
// Tests that care about district derivation use SeedClientAt instead.
func SeedClient(tb testing.TB, gdb *gorm.DB, id string) {
	tb.Helper()
	SeedClientAt(tb, gdb, id, 52.52, 13.405)
}

// SeedClientAt inserts an enrolled station at the given coordinates.
func SeedClientAt(tb testing.TB, gdb *gorm.DB, id string, lat, lng float64) {
	tb.Helper()

	client := database.Client{
		ID:        id,
		Type:      sensor_readings_collector_pkg.ClientTypeDeviceV1,
		Latitude:  lat,
		Longitude: lng,
	}
	if err := gdb.Create(&client).Error; err != nil {
		tb.Fatalf("seed client %s: %v", id, err)
	}
}

// SeedReading inserts one reading for a station at the given measured-at time.
// measuredAt is normalized to UTC, matching what the ingest controller stores.
func SeedReading(tb testing.TB, gdb *gorm.DB, clientID string, pm25, pm10 float64, measuredAt time.Time) {
	tb.Helper()

	reading := database.Reading{
		ClientID:   clientID,
		PM25:       pm25,
		PM10:       pm10,
		MeasuredAt: measuredAt.UTC(),
	}
	if err := gdb.Create(&reading).Error; err != nil {
		tb.Fatalf("seed reading for %s at %s: %v", clientID, measuredAt, err)
	}

	// Maintain the latest-reading cache exactly as ingest does, so seeded data
	// shows up in the fleet snapshot/stats queries that read from it.
	latest := database.LatestReading{ClientID: clientID, PM25: pm25, PM10: pm10, MeasuredAt: measuredAt.UTC()}
	if err := database.UpsertLatestReading(gdb, latest); err != nil {
		tb.Fatalf("seed latest for %s at %s: %v", clientID, measuredAt, err)
	}
}
