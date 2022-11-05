package item

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/recorderrors"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

		err = db.AutoMigrate(&StockOnHand{})
		if err != nil {
			logger.For(ctx).Fatal("Failed to migrate db for type StockOnHand", zap.Error(err))
		}

		r = &repository{db: db}
	}

	return r, nil
}

func (p *repository) GetItem(ctx context.Context, locationID string, itemID string) (*StockOnHand, error) {
	var v StockOnHand
	// TODO: Check the .First query as codegen is not able
	// to elegantly deal with multiple request parameters
	tx := p.db.WithContext(ctx).Model(&StockOnHand{}).First(&v, "location_id = ? AND item_id = ? ", locationID, itemID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) UpdateStockOnHand(ctx context.Context, locationID string, itemID string, delta int) error {
	v := StockOnHand{
		ItemID:      itemID,
		LocationID:  locationID,
		StockOnHand: (delta * -1),
	}
	tx := p.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "item_id"}, {Name: "location_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"stock_on_hand": gorm.Expr("stock_on_hands.stock_on_hand - ?", delta)}),
	}).Create(&v)
	return tx.Error
}
