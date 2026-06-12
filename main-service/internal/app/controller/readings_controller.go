package controller

import (
	"context"
	"errors"
	"go-service-skeleton/internal/app/database"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	sensor_readings_collector_pkg "go-service-skeleton/pkg/sensor-readings-collector"
)

func CreateReading(ctx context.Context, req *sensor_readings_collector_pkg.CreateReadingRequest) (*sensor_readings_collector_pkg.CreateReadingResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)
	var client database.Client

	err := gdb.Select("id").First(&client, "id = ?", req.ClientID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewResponseError(response.ErrorCodeNotFound, response.ErrorMessageNotFound)
		}
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	reading := database.Reading{
		ClientID:   req.ClientID,
		PM25:       *req.Payload.PM25,
		PM10:       *req.Payload.PM10,
		MeasuredAt: req.Payload.MeasuredAt.UTC(),
	}
	result := gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "client_id"}, {Name: "measured_at"}},
		DoNothing: true,
	}).Create(&reading)

	if result.Error != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	return &sensor_readings_collector_pkg.CreateReadingResponse{
		Data: sensor_readings_collector_pkg.Reading{
			ClientID:   reading.ClientID,
			PM25:       reading.PM25,
			PM10:       reading.PM10,
			MeasuredAt: reading.MeasuredAt,
		},
		Created: result.RowsAffected > 0,
	}, nil

}
