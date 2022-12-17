package repository

import (
	"context"

	"github.com/jdotw/warehouse/services/stock-go/pkg/model"
)

type Repository interface {
	GetCategories(ctx context.Context) (*[]model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
}
