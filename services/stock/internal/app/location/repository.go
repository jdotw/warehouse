package location

import (
	"context"
	_ "embed"
)

type Repository interface {
	GetLocations(ctx context.Context) (*[]Location, error)
	CreateLocation(ctx context.Context, location *Location) (*Location, error)
	DeleteLocation(ctx context.Context) error
	GetLocation(ctx context.Context, locationID string) (*Location, error)
	UpdateLocation(ctx context.Context, location *Location) (*Location, error)
}
