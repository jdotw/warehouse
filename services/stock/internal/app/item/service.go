package item

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetItemsInCategory(ctx context.Context, categoryID string) (*[]Item, error)
	CreateItemInCategory(ctx context.Context, item *Item) (*Item, error)
	DeleteItem(ctx context.Context) error
	GetItem(ctx context.Context, itemID string) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) (*Item, error)
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

func (f *service) GetItemsInCategory(ctx context.Context, categoryID string) (*[]Item, error) {
	v, err := f.repository.GetItemsInCategory(ctx, categoryID)
	return v, err
}

func (f *service) CreateItemInCategory(ctx context.Context, item *Item) (*Item, error) {
	v, err := f.repository.CreateItemInCategory(ctx, item)
	return v, err
}

func (f *service) DeleteItem(ctx context.Context) error {
	v, err := f.repository.DeleteItem(ctx)
	return v, err
}

func (f *service) GetItem(ctx context.Context, itemID string) (*Item, error) {
	v, err := f.repository.GetItem(ctx, itemID)
	return v, err
}

func (f *service) UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	v, err := f.repository.UpdateItem(ctx, item)
	return v, err
}
