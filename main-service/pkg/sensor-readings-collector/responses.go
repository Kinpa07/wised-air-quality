package sensor_readings_collector_pkg

type Client struct {
	ID        string         `json:"id"`
	Type      ClientTypeEnum `json:"type"`
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
}

type GetClientsResponse struct {
	Data   []Client `json:"data"`
	Cursor struct {
		After  *string `json:"after"`
		Before *string `json:"before"`
	} `json:"cursor"`
}

type CreateClientResponse struct {
	Data Client `json:"data"`
}
