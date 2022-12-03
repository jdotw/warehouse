package main

import (
	"context"
)

type Service interface {
	GetCategories(ctx context.Context) (*[]Category, error)
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) (*Category, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) (Service, error) {
	return &service{r: r}, nil
}

func (s service) GetCategories(ctx context.Context) (*[]Category, error) {
	c, err := s.r.GetCategories(ctx)
	return c, err
}
func (s service) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	c, err := s.r.CreateCategory(ctx, category)
	return c, err
}
func (s service) DeleteCategory(ctx context.Context, categoryID string) error {
	err := s.r.DeleteCategory(ctx, categoryID)
	return err
}
func (s service) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	c, err := s.r.GetCategory(ctx, categoryID)
	return c, err
}
func (s service) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	c, err := s.r.UpdateCategory(ctx, category)
	return c, err
}
