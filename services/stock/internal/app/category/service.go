package category

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetCategories(ctx context.Context) (*[]Category, error)
	CreateCategory(ctx context.Context, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	GetCategory(ctx context.Context, categoryID string) (*Category, error)
	UpdateCategory(ctx context.Context, category *Category) (*Category, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository, logger log.Factory, tracer opentracing.Tracer) Service {
	var svc Service
	{
		svc = &service{
			repository: repository,
		}
	}
	return svc
}

func (f *service) GetCategories(ctx context.Context) (*[]Category, error) {
	v, err := f.repository.GetCategories(ctx)
	return v, err
}

func (f *service) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	v, err := f.repository.CreateCategory(ctx, category)
	return v, err
}

func (f *service) DeleteCategory(ctx context.Context, categoryID string) error {
	v, err := f.repository.DeleteCategory(ctx, categoryID)
	return v, err
}

func (f *service) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	v, err := f.repository.GetCategory(ctx, categoryID)
	return v, err
}

func (f *service) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	v, err := f.repository.UpdateCategory(ctx, category)
	return v, err
}
