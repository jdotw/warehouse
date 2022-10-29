package transaction

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

		// DO NOT AutoMigrate Transaction and/or TransactionLineItem here
		// They are already AutoMigrated in Item's repository-gorm

		r = &repository{db: db}
	}

	return r, nil
}

func (p *repository) GetTransactions(ctx context.Context) (*[]Transaction, error) {
	var v []Transaction
	tx := p.db.WithContext(ctx).Model(&Transaction{}).Preload("Items").Find(&v)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}

func (p *repository) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	var tx *gorm.DB
	tx = p.db.WithContext(ctx).Create(transaction)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return transaction, nil
}

func (p *repository) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	var v Transaction
	tx := p.db.WithContext(ctx).Model(&Transaction{}).Preload("Items").First(&v, "id = ? ", transactionID)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, recorderrors.ErrNotFound
	}
	return &v, tx.Error
}
