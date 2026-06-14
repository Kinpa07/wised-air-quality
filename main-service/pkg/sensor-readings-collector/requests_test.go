package sensor_readings_collector_pkg

import (
	"math"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
)

func ptr(f float64) *float64 { return &f }

func Test_CreateReadingRequest_Validation(t *testing.T) {
	validate := validator.New()
	// validator.Struct panics on an unregistered custom tag, so wire not_future
	// onto this local instance (init() does it for the router's validator).
	validate.RegisterValidation("not_future", notFuture)
	now := time.Now().UTC()

	cases := []struct {
		name    string
		pm25    *float64
		pm10    *float64
		ts      time.Time
		wantErr bool
	}{
		{"valid reading", ptr(14.2), ptr(23.4), now, false},
		{"zero is a valid reading", ptr(0), ptr(0), now, false},
		{"missing pm2_5", nil, ptr(23.4), now, true},
		{"missing pm10", ptr(14.2), nil, now, true},
		{"negative pm2_5", ptr(-1), ptr(23.4), now, true},
		{"over-ceiling pm2_5", ptr(1500), ptr(23.4), now, true},
		{"over-ceiling pm10", ptr(14.2), ptr(2000), now, true},
		{"non-finite NaN", ptr(math.NaN()), ptr(23.4), now, true},
		{"non-finite +Inf", ptr(math.Inf(1)), ptr(23.4), now, true},
		{"missing timestamp", ptr(14.2), ptr(23.4), time.Time{}, true},
		{"slight clock skew is tolerated", ptr(14.2), ptr(23.4), now.Add(2 * time.Minute), false},
		{"far-future timestamp", ptr(14.2), ptr(23.4), now.Add(24 * time.Hour), true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateReadingRequest{ClientID: "11111111-1111-1111-1111-111111111111"}
			req.Payload.PM25 = tt.pm25
			req.Payload.PM10 = tt.pm10
			req.Payload.MeasuredAt = tt.ts

			err := validate.Struct(req)
			if tt.wantErr && err == nil {
				t.Fatalf("want validation error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}
