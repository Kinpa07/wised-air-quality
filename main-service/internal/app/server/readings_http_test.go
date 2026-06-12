package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"go-service-skeleton/internal/app/database"
	"go-service-skeleton/internal/testutil"
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"

	"github.com/SintroSecurity/go-libraries/logger"
	"github.com/SintroSecurity/go-libraries/metrics"
	"github.com/SintroSecurity/go-libraries/router"
)

const (
	serviceName   = "sensor-readings-collector"
	knownClient   = "11111111-1111-1111-1111-111111111111"
	unknownClient = "99999999-9999-9999-9999-999999999999"
)

// body builds a well-formed reading payload for the given timestamp.
func body(ts string) string {
	return fmt.Sprintf(`{"pm2_5":14.2,"pm10":23.4,"timestamp":%q}`, ts)
}

// sharedHandler is the real router stack, built once for the whole package.
// go-libraries' Prometheus middleware registers its metric vecs against a
// package-level "registered once" guard, so a second router build in the same
// process gets nil gauges and panics — which matches prod (one router per
// process). TestMain honours that by constructing a single router and sharing
// it; inserting cases just use distinct timestamps so they don't collide.
var sharedHandler http.Handler

func TestMain(m *testing.M) {
	gdb, err := testutil.OpenMemDB()

	if err != nil {
		panic(fmt.Sprintf("open mem db: %v", err))
	}

	seed := database.Client{
		ID:        knownClient,
		Type:      sensor_readings_collector_pkg.ClientTypeDeviceV1,
		Latitude:  52.52,
		Longitude: 13.405,
	}

	if err := gdb.Create(&seed).Error; err != nil {
		panic(fmt.Sprintf("seed client: %v", err))
	}

	// The validator middleware reads a logger out of the ctx it captured at
	// registration time, so the ctx handed to InitRoutingTable must carry one —
	// this mirrors StartService, which sets the logger before building routes.
	log := logger.New(context.Background(), &logger.Config{Level: "error"})
	ctx := logger.NewContextWithLogger(context.Background(), log)

	rootMW := CreateRootMiddleware(&Config{Database: gdb})
	r := NewRouter(serviceName, metrics.NewMetrics(), &router.Config{ListenURL: ":0"}, &logger.Config{Level: "error"}, rootMW)
	InitRoutingTable(ctx, r)
	sharedHandler = r

	code := m.Run()
	os.Exit(code)
}

// post fires one request at the {client_id}:read route and returns the status.
func post(t *testing.T, clientID, reqBody string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/clients/%s:read", clientID), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	sharedHandler.ServeHTTP(rec, req)
	return rec.Code
}

func Test_CreateReading_HTTP(t *testing.T) {
	cases := []struct {
		name       string
		clientID   string
		body       string
		wantStatus int
	}{
		// 201 proves both that `{client_id}:read` actually routes through chi and
		// that the full bind→validate→controller→upsert path succeeds.
		{"valid fresh reading", knownClient, body("2026-06-12T09:00:00Z"), http.StatusCreated},

		// out-of-bounds / missing fail the struct-tag validator → 400. That they
		// are NOT 404 also proves validation runs before the client lookup.
		{"missing pm2_5", knownClient, `{"pm10":23.4,"timestamp":"2026-06-12T09:00:00Z"}`, http.StatusBadRequest},
		{"negative pm2_5", knownClient, `{"pm2_5":-3,"pm10":23.4,"timestamp":"2026-06-12T09:00:00Z"}`, http.StatusBadRequest},
		{"over-ceiling pm2_5", knownClient, `{"pm2_5":5000,"pm10":23.4,"timestamp":"2026-06-12T09:00:00Z"}`, http.StatusBadRequest},

		// non-finite and malformed bodies can't be bound → httpin error handler → 422.
		{"non-finite NaN", knownClient, `{"pm2_5":NaN,"pm10":23.4,"timestamp":"2026-06-12T09:00:00Z"}`, http.StatusUnprocessableEntity},
		{"malformed json", knownClient, `{"pm2_5":`, http.StatusUnprocessableEntity},

		// a well-formed reading for a client that was never enrolled → 404.
		{"unknown client", unknownClient, body("2026-06-12T09:00:00Z"), http.StatusNotFound},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := post(t, tt.clientID, tt.body)
			if got != tt.wantStatus {
				t.Fatalf("want %d, got %d", tt.wantStatus, got)
			}
		})
	}
}

// Idempotency on the wire: the first delivery is 201 Created, a re-delivery of
// the same (client_id, timestamp) is skipped by ON CONFLICT DO NOTHING and
// reported as 200 OK — a no-op success, never a 409/duplicate-row error.
// Uses its own timestamp so it doesn't collide with the table test's inserts.
func Test_CreateReading_HTTP_DuplicateIsOk(t *testing.T) {
	dup := body("2026-06-12T10:00:00Z")

	if got := post(t, knownClient, dup); got != http.StatusCreated {
		t.Fatalf("first delivery: want 201, got %d", got)
	}
	if got := post(t, knownClient, dup); got != http.StatusOK {
		t.Fatalf("re-delivery: want 200, got %d", got)
	}
}
