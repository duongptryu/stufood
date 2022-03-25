package ordertrackingstore

import (
	"context"
	"foodlive/modules/ordertracking/ordertrackingmodel"
	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

func NewSqlStore(db *gorm.DB) *sqlStore {
	return &sqlStore{
		db: db,
	}
}

type OrderTrackingStore interface {
	CreateOrderTracking(ctx context.Context, data *ordertrackingmodel.OrderTrackingCreate) error
	UpdateOrder(ctx context.Context, id int, data *ordertrackingmodel.OrderTrackingUpdate) error
	FindOrderTracking(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*ordertrackingmodel.OrderTracking, error)
}
