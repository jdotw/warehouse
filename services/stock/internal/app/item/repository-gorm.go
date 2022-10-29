package item

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/recorderrors"
	"github.com/jdotw/stock/internal/app/location"
	"github.com/jdotw/stock/internal/app/transaction"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

type repository struct {
	db     *gorm.DB
	logger log.Factory
	tracer opentracing.Tracer
}

func NewGormRepository(ctx context.Context, connString string, logger log.Factory, tracer opentracing.Tracer) (Repository, error) {
	var r Repository
	{
		db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to open db", zap.Error(err))
		}

		db.Use(gormopentracing.New(gormopentracing.WithTracer(tracer)))

		// NOTE: We AutoMigrate transaction's Transaction and TransactionLineItem
		// structs here, and also location's Location struct, so that the correct
		// ForeignKey relationship is established between them

		err = db.AutoMigrate(&Item{}, &transaction.Transaction{}, &transaction.TransactionLineItem{}, &location.Location{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type Item", zap.Error(err))
		}

		r = &repository{db: db, logger: logger, tracer: tracer}
	}

	return r, nil
}

func (p *repository) GetItemsInCategory(ctx context.Context, categoryID string) (*[]Item, error) {
	var v []Item
	p.logger.For(ctx).Info("Querying for items", zap.String("category_id", categoryID))
	tx := p.db.WithContext(ctx).Find(&v, "category_id = ? ", categoryID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	p.logger.For(ctx).Info("Found items", zap.String("category_id", categoryID), zap.Int("count", len(v)))
	return &v, tx.Error
}

func (p *repository) CreateItemInCategory(ctx context.Context, item *Item) (*Item, error) {
	var tx *gorm.DB
	tx = p.db.WithContext(ctx).Create(item)
	if tx.Error != nil {
		p.logger.For(ctx).Fatal("CreateItemInCategory Failed", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return item, nil
}

func (p *repository) DeleteItem(ctx context.Context, itemID string) error {
	tx := p.db.WithContext(ctx).Delete(&Item{}, "id = ? ", itemID)
	if tx.RowsAffected == 0 {
		return recorderrors.ErrNotFound
	}
	return tx.Error
}

func (p *repository) GetItem(ctx context.Context, itemID string) (*Item, error) {
	var v Item
	tx := p.db.WithContext(ctx).Model(&Item{}).First(&v, "id = ? ", itemID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateItem(ctx context.Context, item *Item) (*Item, error) {
	var v Item
	tx := p.db.WithContext(ctx).Model(&Item{}).Where("id = ?", item.ID).UpdateColumns(item)
	if tx.RowsAffected == 0 {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
