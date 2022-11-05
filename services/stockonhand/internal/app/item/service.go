package item

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetItem(ctx context.Context, locationID string, itemID string) (*StockOnHand, error)
	UpdateStockOnHand(ctx context.Context, itemID string, delta int) error
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

func (f *service) GetItem(ctx context.Context, locationID string, itemID string) (*StockOnHand, error) {
	v, err := f.repository.GetItem(ctx, locationID, itemID)
	return v, err
}

func (f *service) UpdateStockOnHand(ctx context.Context, locationID string, itemID string, delta int) error {
	err := f.repository.UpdateStockOnHand(ctx, locationID, itemID, delta)
	return err
}
