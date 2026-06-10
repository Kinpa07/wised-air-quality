package db

import (
	"context"
	"gorm.io/gorm"
)

type ctxKeyDatabaseType uint

const ctxKeyDatabase ctxKeyDatabaseType = iota + 1

func NewContextWithDatabase(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyDatabase, db)
}

func GetDatabaseFromContext(ctx context.Context) *gorm.DB {
	if client, err := ctx.Value(ctxKeyDatabase).(*gorm.DB); err != true {
		return nil
	} else {
		return client
	}
}
