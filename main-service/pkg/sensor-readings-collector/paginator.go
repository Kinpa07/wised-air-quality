package sensor_readings_collector_pkg

import (
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

func CreatePaginator(
	cursor paginator.Cursor,
	order *paginator.Order,
	orderFields []string,
	limit *int,
) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Keys:  orderFields,
			Limit: 50,
			Order: paginator.ASC,
		},
	}
	if limit != nil {
		opts = append(opts, paginator.WithLimit(*limit))
	}
	if order != nil {
		opts = append(opts, paginator.WithOrder(*order))
	}
	if cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*cursor.After))
	}
	if cursor.Before != nil {
		opts = append(opts, paginator.WithBefore(*cursor.Before))
	}
	return paginator.New(opts...)
}
