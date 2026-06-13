// Package geo derives a Berlin district from a coordinate via point-in-polygon
// against a static GeoJSON of the 12 Bezirke, loaded once at startup.
package geo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

type ctxKey struct{}

func NewContext(ctx context.Context, idx *DistrictIndex) context.Context {
	return context.WithValue(ctx, ctxKey{}, idx)
}

func FromContext(ctx context.Context) *DistrictIndex {
	idx, _ := ctx.Value(ctxKey{}).(*DistrictIndex)
	return idx
}

type District struct {
	Name    string
	Polygon orb.MultiPolygon
	Bound   orb.Bound
}

type DistrictIndex struct {
	districts []District
}

func (d *DistrictIndex) DistrictFor(lat float64, lng float64) string {
	point := orb.Point{lng, lat}
	for _, district := range d.districts {
		if !district.Bound.Contains(point) {
			continue
		}

		if planar.MultiPolygonContains(district.Polygon, point) {
			return district.Name
		}
	}
	return ""
}

// coordinates decodes straight into orb.MultiPolygon: that type is already the
// nested []Polygon -> []Ring -> []Point shape GeoJSON uses.
type geoJSON struct {
	Features []struct {
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
		Geometry struct {
			Type        string           `json:"type"`
			Coordinates orb.MultiPolygon `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func Load(path string) (*DistrictIndex, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read districts geojson: %w", err)
	}

	var fc geoJSON
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, fmt.Errorf("parse districts geojson: %w", err)
	}

	districts := make([]District, 0, len(fc.Features))
	for _, f := range fc.Features {
		if len(f.Geometry.Coordinates) == 0 {
			continue
		}
		districts = append(districts, District{
			Name:    f.Properties.Name,
			Polygon: f.Geometry.Coordinates,
			Bound:   f.Geometry.Coordinates.Bound(),
		})
	}

	if len(districts) == 0 {
		return nil, fmt.Errorf("districts geojson %q has no usable features", path)
	}

	return &DistrictIndex{districts: districts}, nil
}
