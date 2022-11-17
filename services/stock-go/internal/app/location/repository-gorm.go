package location

import (
	"context"
	_ "embed"

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

		// DO NOT AutoMigrate Location here
		// It is already AutoMigrated in Item's repository-gorm

		r = &repository{db: db}
	}

	return r, nil
}

func (p *repository) GetLocations(ctx context.Context) (*[]Location, error) {
	var v []Location
	tx := p.db.WithContext(ctx).Model(&Location{}).Find(&v)
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
	tx := p.db.WithContext(ctx).Delete(&Location{}, "id = ? ", locationID)
	if tx.RowsAffected == 0 {
		return recorderrors.ErrNotFound
	}
	return tx.Error
}

func (p *repository) GetLocation(ctx context.Context, locationID string) (*Location, error) {
	var v Location
	tx := p.db.WithContext(ctx).Model(&Location{}).First(&v, "id = ? ", locationID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateLocation(ctx context.Context, location *Location) (*Location, error) {
	var v Location
	tx := p.db.WithContext(ctx).Model(&Location{}).Where("id = ?", location.ID).UpdateColumns(location)
	if tx.RowsAffected == 0 {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
