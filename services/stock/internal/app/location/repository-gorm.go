package location

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/recorderrors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

type repository struct {
	db *gorm.DB
}

func NewGormRepository(ctx context.Context, connString string, logger log.Factory, tracer opentracing.Tracer) (Repository, error) {
	var r Repository
	{
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to open db", zap.Error(err))
		}

		db.Use(gormopentracing.New(gormopentracing.WithTracer(tracer)))

		// TODO: Ensure these migrations are correct
		// The OpenAPI Spec used to generate this code often uses
		// results in AutoMigrate statements being generated for
		// request/response body objects instead of actual data models

		err = db.AutoMigrate(&Location{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type Location", zap.Error(err))
		}

		err = db.AutoMigrate(&[]Location{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type []Location", zap.Error(err))
		}

		r = &repository{db: db}
	}

	return r, nil
}

func (p *repository) GetLocations(ctx context.Context) (*[]Location, error) {
	var v []Location
	// TODO: Check the .First query as codegen is not able
	// to elegantly deal with multiple request parameters
	tx := p.db.WithContext(ctx).Model(&[]Location{}).First(&v, "")
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) CreateLocation(ctx context.Context, location *Location) (*Location, error) {
	var tx *gorm.DB
	tx = p.db.WithContext(ctx).Create(location)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return location, nil
}

func (p *repository) DeleteLocation(ctx context.Context, locationID string) error {
	// TODO: Unable to generate code for this Operation
	return nil, errors.New("Not Implemented")
}

func (p *repository) GetLocation(ctx context.Context, locationID string) (*Location, error) {
	var v Location
	// TODO: Check the .First query as codegen is not able
	// to elegantly deal with multiple request parameters
	tx := p.db.WithContext(ctx).Model(&Location{}).First(&v, "location_id = ? ", locationID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateLocation(ctx context.Context, location *Location) (*Location, error) {
	var v Location
	// TODO: Check that the .Where query is appropriate
	tx := p.db.WithContext(ctx).Model(&Location{}).Where("id = ?", location.ID).UpdateColumns(location)
	if tx.RowsAffected == 0 {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
