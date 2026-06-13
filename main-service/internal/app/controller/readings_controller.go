package controller

import (
	"context"
	"errors"
	"go-service-skeleton/internal/app/database"

	"github.com/SintroSecurity/go-libraries/db"
	"github.com/SintroSecurity/go-libraries/router/response"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
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
	latest := database.LatestReading{
		ClientID:   reading.ClientID,
		PM25:       reading.PM25,
		PM10:       reading.PM10,
		MeasuredAt: reading.MeasuredAt,
	}

	var created bool
	err = gdb.Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "client_id"}, {Name: "measured_at"}},
			DoNothing: true,
		}).Create(&reading)
		if result.Error != nil {
			return result.Error
		}
		created = result.RowsAffected > 0

		return database.UpsertLatestReading(tx, latest)
	})
	if err != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	return &sensor_readings_collector_pkg.CreateReadingResponse{
		Data: sensor_readings_collector_pkg.Reading{
			ClientID:   reading.ClientID,
			PM25:       reading.PM25,
			PM10:       reading.PM10,
			MeasuredAt: reading.MeasuredAt,
		},
		Created: created,
	}, nil

}

func GetReading(ctx context.Context, req *sensor_readings_collector_pkg.GetReadingsRequest) (*sensor_readings_collector_pkg.GetReadingsResponse, *response.Error) {
	gdb := db.GetDatabaseFromContext(ctx)

	var client database.Client
	err := gdb.Select("id").First(&client, "id = ?", req.ClientID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewResponseError(response.ErrorCodeNotFound, response.ErrorMessageNotFound)
		}
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	query := gdb.Where("client_id = ?", req.ClientID)

	if req.From != nil && req.To != nil {
		query = query.Where("measured_at BETWEEN ? AND ?", *req.From, *req.To)
	} else if req.Since != nil {
		query = query.Where("measured_at >= ?", *req.Since)
	}

	order := paginator.DESC
	p := sensor_readings_collector_pkg.CreatePaginator(
		paginator.Cursor{After: req.After, Before: req.Before},
		&order,
		[]string{"MeasuredAt"},
		req.Limit,
	)

	var readings []database.Reading

	_, cursor, err := p.Paginate(query, &readings)
	if err != nil {
		return nil, response.NewResponseError(response.ErrorCodeInternal, response.ErrorMessageInternalServerError)
	}

	data := make([]sensor_readings_collector_pkg.Reading, len(readings))

	for i, u := range readings {
		data[i] = sensor_readings_collector_pkg.Reading{
			ClientID:   u.ClientID,
			PM25:       u.PM25,
			PM10:       u.PM10,
			MeasuredAt: u.MeasuredAt,
		}
	}

	res := &sensor_readings_collector_pkg.GetReadingsResponse{Data: data}
	res.Cursor.After = cursor.After
	res.Cursor.Before = cursor.Before

	return res, nil
}
