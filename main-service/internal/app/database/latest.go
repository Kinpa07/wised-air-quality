package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpsertLatestReading advances a station's cached latest reading, but only when
// the incoming reading is newer than the stored one — so a late older reading
// or a duplicate re-delivery is a no-op. Shared by ingest and the seed helper.
func UpsertLatestReading(db *gorm.DB, r LatestReading) error {
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "client_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"pm2_5":       r.PM25,
			"pm10":        r.PM10,
			"measured_at": r.MeasuredAt,
		}),
		Where: clause.Where{Exprs: []clause.Expression{
			clause.Expr{SQL: "latest_readings.measured_at < ?", Vars: []interface{}{r.MeasuredAt}},
		}},
	}).Create(&r).Error
}
