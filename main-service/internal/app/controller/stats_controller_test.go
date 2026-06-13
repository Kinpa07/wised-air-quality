package controller

import (
	"context"
	"testing"
	"time"

	"go-service-skeleton/internal/app/display"
	"go-service-skeleton/internal/testutil"

	"gorm.io/gorm"
)

const (
	statsA = "aaaaaaaa-0000-0000-0000-000000000001"
	statsB = "bbbbbbbb-0000-0000-0000-000000000002"
	statsC = "cccccccc-0000-0000-0000-000000000003"
)

// statsCtx tailors the windows so the fixture is hand-computable without seeding
// hundreds of rows: connection 60 / interval 30 => expected 2, so 2 readings = ratio
// 1.0 and 1 reading = 0.5. The active window is 15 min, distinct from the 60-min
// connection window, so a reading can count toward the ratio but still be "stale".
func statsCtx(t *testing.T) (context.Context, *gorm.DB) {
	t.Helper()
	ctx, gdb := ctxWithDB(t)

	ctx = display.NewContext(ctx, &display.Config{
		ActiveWindowMinutes:     15,
		ConnectionWindowMinutes: 60,
		ExpectedIntervalMinutes: 30,
		PoorConnectionThreshold: 0.8,
	})

	return ctx, gdb
}

// Mixed fleet, all numbers pre-computed:
//
//	A: 2 readings in the window, newest 5 min ago => active, ratio 1.0, pm 14
//	B: 1 reading 40 min ago        => stale (not active), ratio 0.5 (poor)
//	C: no readings                 => silent, ratio 0.0 (poor)
//
// active=1, avgPM25=14 (only A), poor=2 (B,C), stability=(1.0+0.5+0.0)/3*100=50.
func Test_GetStats_MixedFleet(t *testing.T) {
	ctx, gdb := statsCtx(t)
	now := time.Now().UTC()

	testutil.SeedClient(t, gdb, statsA)
	testutil.SeedReading(t, gdb, statsA, 9, 18, now.Add(-40*time.Minute))
	testutil.SeedReading(t, gdb, statsA, 14, 22, now.Add(-5*time.Minute)) // newest, active

	testutil.SeedClient(t, gdb, statsB)
	testutil.SeedReading(t, gdb, statsB, 30, 50, now.Add(-40*time.Minute)) // stale

	testutil.SeedClient(t, gdb, statsC) // silent

	res, errResp := GetStats(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	s := res.Data

	if s.ActiveSensors != 1 {
		t.Fatalf("active sensors: want 1, got %d", s.ActiveSensors)
	}
	if s.AvgPM25 == nil || *s.AvgPM25 != 14 {
		t.Fatalf("avg pm2_5: want 14, got %v", s.AvgPM25)
	}
	if s.PoorConnection != 2 {
		t.Fatalf("poor connection: want 2, got %d", s.PoorConnection)
	}
	if s.NetworkStability != 50 {
		t.Fatalf("network stability: want 50, got %v", s.NetworkStability)
	}
}

// Empty fleet: no NaN from 0/0 — avg is nil, stability is 0.
func Test_GetStats_EmptyFleet(t *testing.T) {
	ctx, _ := statsCtx(t)

	res, errResp := GetStats(ctx)
	if errResp != nil {
		t.Fatalf("unexpected error: %v", errResp)
	}
	s := res.Data

	if s.ActiveSensors != 0 || s.PoorConnection != 0 {
		t.Fatalf("counts: want 0/0, got %d/%d", s.ActiveSensors, s.PoorConnection)
	}
	if s.AvgPM25 != nil {
		t.Fatalf("avg pm2_5: want nil, got %v", *s.AvgPM25)
	}
	if s.NetworkStability != 0 {
		t.Fatalf("network stability: want 0, got %v", s.NetworkStability)
	}
}
