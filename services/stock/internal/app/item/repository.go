package item

import (
	"context"
	_ "embed"
)

type Repository interface {
	GetItemsInCategory(ctx context.Context, categoryID string) (*[]Item, error)
	CreateItemInCategory(ctx context.Context, item *Item) (*Item, error)
	DeleteItem(ctx context.Context, itemID string) error
	GetItem(ctx context.Context, itemID string) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) (*Item, error)
}
