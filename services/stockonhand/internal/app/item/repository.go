package item

import (
	"context"
	_ "embed"
)

type Repository interface {
	GetItem(ctx context.Context, locationID string, itemID string) (*StockOnHand, error)
	UpdateStockOnHand(ctx context.Context, locationID string, itemID string, delta int) (int, error)
}
