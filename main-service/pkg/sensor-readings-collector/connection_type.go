package sensor_readings_collector_pkg

type ConnectionQuality string

const (
	ConnectionGood ConnectionQuality = "Good"
	ConnectionPoor ConnectionQuality = "Poor"
)

func QualityFor(ratio float64, threshold float64) ConnectionQuality {
	if ratio >= threshold {
		return ConnectionGood
	}
	return ConnectionPoor

}
