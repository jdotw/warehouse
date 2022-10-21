package category

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
	}

	return r, nil
}

func (p *repository) GetCategories(ctx context.Context) (*[]Category, error) {
	var v []Category
	// TODO: Check the .First query as codegen is not able
	// to elegantly deal with multiple request parameters
	tx := p.db.WithContext(ctx).Model(&[]Category{}).First(&v, "")
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
	// TODO: Unable to generate code for this Operation
	return nil
}

func (p *repository) GetCategory(ctx context.Context, categoryID string) (*Category, error) {
	var v Category
	// TODO: Check the .First query as codegen is not able
	// to elegantly deal with multiple request parameters
	tx := p.db.WithContext(ctx).Model(&Category{}).First(&v, "category_id = ? ", categoryID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateCategory(ctx context.Context, category *Category) (*Category, error) {
	var v Category
	// TODO: Check that the .Where query is appropriate
	tx := p.db.WithContext(ctx).Model(&Category{}).Where("id = ?", category.ID).UpdateColumns(category)
	if tx.RowsAffected == 0 {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
