package database

import (
	"fmt"

	"gorm.io/gorm"
)

func MigrateDatabase(gdb *gorm.DB) error {
	if err := Migrate(gdb); err != nil {
		return fmt.Errorf("database migrate error: %w", err)
	}
	return nil
}

func Migrate(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(&Client{}, &Reading{})
	})
}
