package controller

import (
	"testing"
	"time"

	"go-service-skeleton/internal/testutil"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"

	"github.com/SintroSecurity/go-libraries/router/response"
)

// intPtr / timePtr keep the optional-pointer request fields readable in tests.
func intPtr(n int) *int              { return &n }
func timePtr(t time.Time) *time.Time { return &t }

func Test_GetReading_UnknownClient(t *testing.T) {
	ctx, _ := ctxWithDB(t)

	_, errResp := GetReading(ctx, &sensor_readings_collector_pkg.GetReadingsRequest{ClientID: unknownClient})

	if errResp == nil {
		t.Fatalf("want 404 error, got nil")
	}
	if errResp.Code != response.ErrorCodeNotFound {
		t.Fatalf("want %s, got %s", response.ErrorCodeNotFound, errResp.Code)
	}
}

func Test_GetReading_ZeroReadingStation(t *testing.T) {
	ctx, gdb := ctxWithDB(t)
	testutil.SeedClient(t, gdb, knownClient) // enrolled, but no readings

	res, errResp := GetReading(ctx, &sensor_readings_collector_pkg.GetReadingsRequest{ClientID: knownClient})

	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	if len(res.Data) != 0 {
		t.Fatalf("want empty data, got %d rows", len(res.Data))
	}
}

func Test_GetReading_SinceFilter(t *testing.T) {
	ctx, gdb := ctxWithDB(t)
	testutil.SeedClient(t, gdb, knownClient)

	base := time.Date(2026, 6, 13, 8, 0, 0, 0, time.UTC)
	testutil.SeedReading(t, gdb, knownClient, 10, 20, base)                  // excluded
	testutil.SeedReading(t, gdb, knownClient, 11, 21, base.Add(1*time.Hour)) // included (== since)
	testutil.SeedReading(t, gdb, knownClient, 12, 22, base.Add(2*time.Hour)) // included

	since := base.Add(1 * time.Hour)
	res, errResp := GetReading(ctx, &sensor_readings_collector_pkg.GetReadingsRequest{
		ClientID: knownClient,
		Since:    timePtr(since),
	})

	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	if len(res.Data) != 2 {
		t.Fatalf("want 2 rows at/after since, got %d", len(res.Data))
	}
	// DESC order: every returned row must be >= since.
	for _, r := range res.Data {
		if r.MeasuredAt.Before(since) {
			t.Fatalf("row %s is older than since %s", r.MeasuredAt, since)
		}
	}
}

func Test_GetReading_PaginationCursor(t *testing.T) {
	ctx, gdb := ctxWithDB(t)
	testutil.SeedClient(t, gdb, knownClient)

	base := time.Date(2026, 6, 13, 8, 0, 0, 0, time.UTC)
	// 4 readings, one minute apart. DESC => newest (base+3) first.
	for i := 0; i < 4; i++ {
		testutil.SeedReading(t, gdb, knownClient, 10, 20, base.Add(time.Duration(i)*time.Minute))
	}

	// Page 1: limit 2 => the two newest (base+3, base+2).
	page1, errResp := GetReading(ctx, &sensor_readings_collector_pkg.GetReadingsRequest{
		ClientID: knownClient,
		Limit:    intPtr(2),
	})
	if errResp != nil {
		t.Fatalf("page 1 error: %v", errResp)
	}
	if len(page1.Data) != 2 {
		t.Fatalf("page 1: want 2 rows, got %d", len(page1.Data))
	}
	if page1.Cursor.After == nil {
		t.Fatalf("page 1: want a non-nil After cursor (more rows exist)")
	}

	// Page 2: follow the After cursor => the two oldest (base+1, base).
	page2, errResp := GetReading(ctx, &sensor_readings_collector_pkg.GetReadingsRequest{
		ClientID: knownClient,
		Limit:    intPtr(2),
		After:    page1.Cursor.After,
	})
	if errResp != nil {
		t.Fatalf("page 2 error: %v", errResp)
	}
	if len(page2.Data) != 2 {
		t.Fatalf("page 2: want 2 rows, got %d", len(page2.Data))
	}

	// No overlap: every page-2 timestamp must be strictly older than every page-1 one.
	oldestPage1 := page1.Data[len(page1.Data)-1].MeasuredAt
	for _, r := range page2.Data {
		if !r.MeasuredAt.Before(oldestPage1) {
			t.Fatalf("page 2 row %s overlaps page 1 (oldest %s)", r.MeasuredAt, oldestPage1)
		}
	}
}
