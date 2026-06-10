package controller

import (
	"context"
	"errors"
	"go-service-skeleton/internal/app/database"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"gorm.io/gorm"

	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
)

func GetClients(ctx context.Context, req *sensor_readings_collector_pkg.GetClientsRequest) (*sensor_readings_collector_pkg.GetClientsResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)

	p := sensor_readings_collector_pkg.CreatePaginator(
		paginator.Cursor{After: req.After, Before: req.Before},
		nil,
		[]string{"ID"},
		req.Limit,
	)

	var clients []database.Client
	_, cursor, err := p.Paginate(gdb, &clients)
	if err != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	data := make([]sensor_readings_collector_pkg.Client, len(clients))
	for i, u := range clients {
		data[i] = sensor_readings_collector_pkg.Client{
			ID:        u.ID,
			Type:      u.Type.ToEnum(),
			Latitude:  u.Latitude,
			Longitude: u.Longitude,
		}
	}

	res := &sensor_readings_collector_pkg.GetClientsResponse{Data: data}
	res.Cursor.After = cursor.After
	res.Cursor.Before = cursor.Before

	return res, nil
}

func CreateClient(ctx context.Context, req *sensor_readings_collector_pkg.CreateClientRequest) (*sensor_readings_collector_pkg.CreateClientResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)

	client := database.Client{
		ID:        req.Payload.ID,
		Type:      req.Payload.ClientType,
		Latitude:  req.Payload.Latitude,
		Longitude: req.Payload.Longitude,
	}

	if err := gdb.Create(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, response.NewResponseError(response.ErrorCodeConflict, response.ErrorMessageConflict)
		}
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	return &sensor_readings_collector_pkg.CreateClientResponse{
		Data: sensor_readings_collector_pkg.Client{
			ID:        client.ID,
			Type:      client.Type.ToEnum(),
			Latitude:  client.Latitude,
			Longitude: client.Longitude,
		},
	}, nil
}
