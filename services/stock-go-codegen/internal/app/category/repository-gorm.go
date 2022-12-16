package category

import (
	"context"
	_ "embed"
	"time"

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

		maxOpenConn := 100

		sqlDB, err := db.DB()
		sqlDB.SetMaxIdleConns(maxOpenConn)
		sqlDB.SetMaxOpenConns(maxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// TODO: Ensure these migrations are correct
		// The OpenAPI Spec used to generate this code often uses
		// results in AutoMigrate statements being generated for
		// request/response body objects instead of actual data models

		err = db.AutoMigrate(&Category{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type Category", zap.Error(err))
		}

		err = db.AutoMigrate(&[]Category{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type []Category", zap.Error(err))
		}

		r = &repository{db: db}

		// Preheat the DB connections by pinging
		// in parallel up to the maxOpenConn count
		for i := 0; i < maxOpenConn; i++ {
			go sqlDB.Ping()
		}

	}

	return r, nil
}

func (p *repository) GetCategories(ctx context.Context) (*[]Category, error) {
	var v []Category
	tx := p.db.WithContext(ctx).Find(&v)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) CreateCategory(ctx context.Context, category *Category) (*Category, error) {
	var tx *gorm.DB
	tx = p.db.WithContext(ctx).Create(category)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return category, nil
}

func (p *repository) DeleteCategory(ctx context.Context, categoryID string) error {
	tx := p.db.WithContext(ctx).Delete(&Category{}, "id = ? ", categoryID)
	if tx.RowsAffected == 0 {
		return recorderrors.ErrNotFound
	}
	return tx.Error
}

func (p *repository) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	var v Category
	tx := p.db.WithContext(ctx).Model(&Category{}).First(&v, "id = ? ", categoryID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	var v Category
	tx := p.db.WithContext(ctx).Model(&Category{}).Where("id = ?", category.ID).UpdateColumns(category)
	if tx.RowsAffected == 0 {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
