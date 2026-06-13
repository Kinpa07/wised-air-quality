package server

import (
	"go-service-skeleton/internal/app/display"
	"go-service-skeleton/internal/app/geo"

	"gorm.io/gorm"
)

type Config struct {
	Database  *gorm.DB
	Display   *display.Config
	Districts *geo.DistrictIndex
}
