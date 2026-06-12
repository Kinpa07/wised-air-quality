package controller

import (
	"context"
	"testing"
	"time"

	"go-service-skeleton/internal/app/database"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	knownClient   = "11111111-1111-1111-1111-111111111111"
	unknownClient = "99999999-9999-9999-9999-999999999999"
)

func setupTestDB(t *testing.T) (context.Context, *gorm.DB) {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("open in-memory sqlite: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1) // keep the in-memory DB on one connection
	t.Cleanup(func() { _ = sqlDB.Close() })

	if err := database.Migrate(gdb); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	ctx := db.NewContextWithDatabase(context.Background(), gdb)
	return ctx, gdb
}

func seedClient(t *testing.T, gdb *gorm.DB, id string) {
	t.Helper()

	client := database.Client{
		ID:        id,
		Type:      sensor_readings_collector_pkg.ClientTypeDeviceV1,
		Latitude:  52.52,
		Longitude: 13.405,
	}
	if err := gdb.Create(&client).Error; err != nil {
		t.Fatalf("seed client: %v", err)
	}
}

func readingReq(clientID string, pm25, pm10 float64, ts time.Time) *sensor_readings_collector_pkg.CreateReadingRequest {
	r := &sensor_readings_collector_pkg.CreateReadingRequest{ClientID: clientID}
	r.Payload.PM25 = &pm25
	r.Payload.PM10 = &pm10
	r.Payload.MeasuredAt = ts
	return r
}

func Test_CreateReading(t *testing.T) {
	ts := time.Date(2026, 6, 12, 9, 0, 0, 0, time.UTC)

	cases := []struct {
		name        string
		seed        bool
		req         *sensor_readings_collector_pkg.CreateReadingRequest
		callTwice   bool
		wantErrCode response.ErrorCode
		wantCreated bool
		wantCount   int64
	}{
		{
			name:        "valid reading is inserted",
			seed:        true,
			req:         readingReq(knownClient, 14.2, 23.4, ts),
			wantCreated: true,
			wantCount:   1,
		},
		{
			name:        "unknown client returns not found",
			req:         readingReq(unknownClient, 14.2, 23.4, ts),
			wantErrCode: response.ErrorCodeNotFound,
			wantCount:   0,
		},
		{
			name:        "duplicate delivery is a no-op",
			seed:        true,
			req:         readingReq(knownClient, 14.2, 23.4, ts),
			callTwice:   true,
			wantCreated: false,
			wantCount:   1,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctx, gdb := setupTestDB(t)
			if tt.seed {
				seedClient(t, gdb, tt.req.ClientID)
			}

			res, errResp := CreateReading(ctx, tt.req)
			if tt.callTwice {
				if errResp != nil {
					t.Fatalf("first call errored: %v", errResp)
				}
				res, errResp = CreateReading(ctx, tt.req)
			}

			if tt.wantErrCode != "" {
				if errResp == nil {
					t.Fatalf("want error %s, got nil", tt.wantErrCode)
				}
				if errResp.Code != tt.wantErrCode {
					t.Fatalf("want error code %s, got %s", tt.wantErrCode, errResp.Code)
				}
			} else {
				if errResp != nil {
					t.Fatalf("unexpected error: %v", errResp)
				}
				if res.Created != tt.wantCreated {
					t.Fatalf("want Created=%v, got %v", tt.wantCreated, res.Created)
				}
			}

			var count int64
			gdb.Model(&database.Reading{}).Count(&count)
			if count != tt.wantCount {
				t.Fatalf("want %d rows, got %d", tt.wantCount, count)
			}
		})
	}
}

func Test_CreateReading_NormalizesMeasuredAtToUTC(t *testing.T) {
	ctx, gdb := setupTestDB(t)
	seedClient(t, gdb, knownClient)

	berlin := time.FixedZone("CEST", 2*60*60)
	local := time.Date(2026, 6, 12, 11, 0, 0, 0, berlin)
	wantUTC := time.Date(2026, 6, 12, 9, 0, 0, 0, time.UTC)

	res, errResp := CreateReading(ctx, readingReq(knownClient, 14.2, 23.4, local))
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	if !res.Data.MeasuredAt.Equal(wantUTC) {
		t.Fatalf("want %s, got %s", wantUTC, res.Data.MeasuredAt)
	}
	if res.Data.MeasuredAt.Location() != time.UTC {
		t.Fatalf("measured_at not in UTC: %s", res.Data.MeasuredAt.Location())
	}

	var stored database.Reading
	if err := gdb.First(&stored).Error; err != nil {
		t.Fatalf("load stored reading: %v", err)
	}
	if !stored.MeasuredAt.Equal(wantUTC) {
		t.Fatalf("stored measured_at want %s, got %s", wantUTC, stored.MeasuredAt)
	}
}
