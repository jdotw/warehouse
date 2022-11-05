package transaction

import (
	"context"
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/recorderrors"
	"github.com/opentracing/opentracing-go"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
)

type repository struct {
	db *gorm.DB
}

func (t *Transaction) AfterSave(tx *gorm.DB) (err error) {
	var wg sync.WaitGroup
	for _, item := range t.Items {
		wg.Add(1)
		event := TransactionLineItemCreatedEvent{
			LineItem: item,
		}
		json, err := json.Marshal(event)
		if err != nil {
			return err
			t.logger.For(t.ctx).Error("Failed to marshall transaction line item", zap.Error(err))
		} else {
			record := &kgo.Record{Topic: "warehouse.stock.transaction.line_item.created", Value: json}
			t.kafka.Produce(t.ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					t.logger.For(t.ctx).Error("Failed to produce event", zap.Error(err))
				}
			})
		}
	}
	wg.Wait()
	return
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
