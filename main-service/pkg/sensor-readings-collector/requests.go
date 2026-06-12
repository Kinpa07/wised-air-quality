package sensor_readings_collector_pkg

import "time"

type CreateClientRequest struct {
	Payload struct {
		ID         string     `json:"id" validate:"required,uuid4"`
		ClientType ClientType `json:"client_type" validate:"client_type"`
		Latitude   float64    `json:"latitude" validate:"required,gte=-90,lte=90"`
		Longitude  float64    `json:"longitude" validate:"required,gte=-180,lte=180"`
	} `in:"body=json"`
}

type GetClientsRequest struct {
	After  *string `in:"query=after"`
	Before *string `in:"query=before"`
	Limit  *int    `in:"query=limit"`
}

type CreateReadingRequest struct {
	ClientID string `in:"path=client_id"`
	Payload  struct {
		PM25       *float64  `json:"pm2_5" validate:"required,gte=0,lte=1000"`
		PM10       *float64  `json:"pm10" validate:"required,gte=0,lte=1000"`
		MeasuredAt time.Time `json:"timestamp" validate:"required"`
	} `in:"body=json"`
}
