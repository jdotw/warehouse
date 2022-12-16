package main

import (
	"context"
)

type Repository interface {
	GetCategories(ctx context.Context) (*[]Category, error)
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) (*Category, error)
}
