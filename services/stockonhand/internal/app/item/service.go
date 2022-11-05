package item

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetItem(ctx context.Context, itemID string) (*ItemStockOnHand, error)
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

func (f *service) GetItem(ctx context.Context, itemID string) (*ItemStockOnHand, error) {
	v, err := f.repository.GetItem(ctx, itemID)
	return v, err
}
