package sensor_readings_collector_pkg

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
