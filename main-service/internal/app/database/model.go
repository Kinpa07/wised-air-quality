package database

import (
	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID        string `gorm:"size:36"`
	Type      sensor_readings_collector_pkg.ClientType
	Latitude  float64
	Longitude float64
	Model
}

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Reading struct {
	ClientID   string  `gorm:"size:36;uniqueIndex:idx_reading_dedup;index:idx_reading_latest_lookup"`
	PM25       float64 `gorm:"column:pm2_5"`
	PM10       float64
	MeasuredAt time.Time `gorm:"uniqueIndex:idx_reading_dedup;index:idx_reading_latest_lookup,sort:desc"`
	CreatedAt  time.Time
}
