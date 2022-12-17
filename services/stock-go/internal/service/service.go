package service

import (
	"context"

	"github.com/jdotw/warehouse/services/stock-go/internal/repository"
	"github.com/jdotw/warehouse/services/stock-go/pkg/model"
)

type Service interface {
	GetCategories(ctx context.Context) (*[]model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
}

type service struct {
	r repository.Repository
}

func NewService(r repository.Repository) (Service, error) {
	return &service{r: r}, nil
}

func (s service) GetCategories(ctx context.Context) (*[]model.Category, error) {
	c, err := s.r.GetCategories(ctx)
	return c, err
}
func (s service) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	c, err := s.r.CreateCategory(ctx, category)
	return c, err
}
func (s service) DeleteCategory(ctx context.Context, categoryID string) error {
	err := s.r.DeleteCategory(ctx, categoryID)
	return err
}
func (s service) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	c, err := s.r.GetCategory(ctx, categoryID)
	return c, err
}
func (s service) UpdateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	c, err := s.r.UpdateCategory(ctx, category)
	return c, err
}
