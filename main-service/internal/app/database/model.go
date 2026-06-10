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
