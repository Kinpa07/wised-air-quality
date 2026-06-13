package geo

import (
	"sync"
	"testing"
)

const districtsPath = "../../../assets/berlin-districts.geojson"

func Test_Load_Error(t *testing.T) {
	if _, err := Load("does-not-exist.geojson"); err == nil {
		t.Fatalf("want error for missing file, got nil")
	}
}

func Test_DistrictFor(t *testing.T) {
	idx, err := Load(districtsPath)
	if err != nil {
		t.Fatalf("load districts: %v", err)
	}

	cases := []struct {
		name     string
		lat, lng float64
		want     string
	}{
		{"Brandenburg Gate", 52.5163, 13.3777, "Mitte"},
		{"Alexanderplatz", 52.5219, 13.4132, "Mitte"},
		{"Köpenick", 52.4445, 13.5743, "Treptow-Köpenick"},
		{"Charlottenburg Palace", 52.5208, 13.2957, "Charlottenburg-Wilmersdorf"},
		{"Munich (outside Berlin)", 48.1374, 11.5755, ""},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.DistrictFor(tt.lat, tt.lng)
			if got != tt.want {
				t.Fatalf("DistrictFor(%v, %v) = %q, want %q", tt.lat, tt.lng, got, tt.want)
			}
		})
	}
}

func Test_DistrictForClient_MatchesDistrictFor(t *testing.T) {
	idx, err := Load(districtsPath)
	if err != nil {
		t.Fatalf("load districts: %v", err)
	}

	// First call (miss, fills cache) and second call (hit) both agree with the
	// uncached DistrictFor.
	want := idx.DistrictFor(52.5163, 13.3777) // Mitte
	for i := 0; i < 2; i++ {
		if got := idx.DistrictForClient("client-1", 52.5163, 13.3777); got != want {
			t.Fatalf("call %d: DistrictForClient = %q, want %q", i, got, want)
		}
	}
}

// Run with -race: many goroutines hammering the shared cache concurrently must
// not trip the race detector or panic on concurrent map access. Mixes repeats of
// one id (read contention on a hit) with distinct ids (write contention on misses).
func Test_DistrictForClient_Concurrent(t *testing.T) {
	idx, err := Load(districtsPath)
	if err != nil {
		t.Fatalf("load districts: %v", err)
	}

	ids := []string{"a", "a", "a", "b", "c", "d"}
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		for _, id := range ids {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				idx.DistrictForClient(id, 52.5163, 13.3777)
			}(id)
		}
	}
	wg.Wait()
}
