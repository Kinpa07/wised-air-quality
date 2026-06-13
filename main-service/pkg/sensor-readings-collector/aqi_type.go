package sensor_readings_collector_pkg

type AQIBand string

const (
	AQIBandGood      AQIBand = "Good"
	AQIBandModerate  AQIBand = "Moderate"
	AQIBandElevated  AQIBand = "Elevated"
	AQIBandUnhealthy AQIBand = "Unhealthy"
)

func BandFor(pm *float64) *AQIBand {
	if pm == nil {
		return nil
	}

	v := *pm
	var band AQIBand
	switch {
	case v < 12.0:
		band = AQIBandGood
	case v < 35.0:
		band = AQIBandModerate
	case v < 55.0:
		band = AQIBandElevated
	default:
		band = AQIBandUnhealthy
	}
	return &band
}
