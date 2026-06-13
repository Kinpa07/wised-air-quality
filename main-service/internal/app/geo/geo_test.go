package geo

import "testing"

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
