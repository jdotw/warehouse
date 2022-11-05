package item

import (
	"context"
	_ "embed"
)

type Repository interface {
	GetItem(ctx context.Context, itemID string) (*ItemStockOnHand, error)
	UpdateStockOnHand(ctx context.Context, itemID string, delta int) error
}
