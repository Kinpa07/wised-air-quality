package server

import (
	"gorm.io/gorm"
)

type Config struct {
	Database *gorm.DB
}
