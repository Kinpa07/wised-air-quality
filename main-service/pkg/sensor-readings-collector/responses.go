package sensor_readings_collector_pkg

import "time"

type Client struct {
	ID        string         `json:"id"`
	Type      ClientTypeEnum `json:"type"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
}

type Reading struct {
	ClientID   string    `json:"client_id"`
	PM25       float64   `json:"pm2_5"`
	PM10       float64   `json:"pm10"`
	MeasuredAt time.Time `json:"timestamp"`
}
type Station struct {
	ID         string            `json:"id"`
	Latitude   float64           `json:"lat"`
	Longitude  float64           `json:"lng"`
	PM25       *float64          `json:"pm2_5"`
	PM10       *float64          `json:"pm10"`
	MeasuredAt *time.Time        `json:"measured_at"`
	Band       *AQIBand          `json:"band"`
	District   string            `json:"district"`
	Stability  float64           `json:"stability"`
	Connection ConnectionQuality `json:"connection"`
}

type Cursor struct {
	After  *string `json:"after"`
	Before *string `json:"before"`
}

type GetClientsResponse struct {
	Data   []Client `json:"data"`
	Cursor Cursor   `json:"cursor"`
}

type CreateClientResponse struct {
	Data Client `json:"data"`
}

type CreateReadingResponse struct {
	Data    Reading `json:"data"`
	Created bool    `json:"-"`
}

type GetReadingsResponse struct {
	Data   []Reading `json:"data"`
	Cursor Cursor    `json:"cursor"`
}

type GetStationsResponse struct {
	Data []Station `json:"data"`
}
