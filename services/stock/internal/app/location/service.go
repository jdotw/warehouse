package location

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetLocations(ctx context.Context) (*[]Location, error)
	CreateLocation(ctx context.Context, location *Location) (*Location, error)
	DeleteLocation(ctx context.Context) error
	GetLocation(ctx context.Context, locationID string) (*Location, error)
	UpdateLocation(ctx context.Context, location *Location) (*Location, error)
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

func (f *service) GetLocations(ctx context.Context) (*[]Location, error) {
	v, err := f.repository.GetLocations(ctx)
	return v, err
}

func (f *service) CreateLocation(ctx context.Context, location *Location) (*Location, error) {
	v, err := f.repository.CreateLocation(ctx, location)
	return v, err
}

func (f *service) DeleteLocation(ctx context.Context) error {
	v, err := f.repository.DeleteLocation(ctx)
	return v, err
}

func (f *service) GetLocation(ctx context.Context, locationID string) (*Location, error) {
	v, err := f.repository.GetLocation(ctx, locationID)
	return v, err
}

func (f *service) UpdateLocation(ctx context.Context, location *Location) (*Location, error) {
	v, err := f.repository.UpdateLocation(ctx, location)
	return v, err
}
