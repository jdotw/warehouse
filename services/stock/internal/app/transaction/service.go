package transaction

import (
	"context"
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Service interface {
	GetTransactions(ctx context.Context) (*[]Transaction, error)
	CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*Transaction, error)
}

type service struct {
	repository Repository
	kafka      *kgo.Client
	logger     log.Factory
}

func NewService(repository Repository, kafka *kgo.Client, logger log.Factory, tracer opentracing.Tracer) Service {
	var svc Service
	{
		svc = &service{
			repository: repository,
			kafka:      kafka,
			logger:     logger,
		}
	}
	return svc
}

func (f *service) GetTransactions(ctx context.Context) (*[]Transaction, error) {
	v, err := f.repository.GetTransactions(ctx)
	return v, err
}

func (f *service) CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error) {
	v, err := f.repository.CreateTransaction(ctx, transaction)

	// 1.) Producing a message
	// All record production goes through Produce, and the callback can be used
	// to allow for synchronous or asynchronous production.
	var wg sync.WaitGroup
	for _, item := range transaction.Items {
		wg.Add(1)
		event := TransactionLineItemCreatedEvent{
			LineItem: item,
		}
		json, err := json.Marshal(event)
		if err != nil {
			f.logger.For(ctx).Error("Failed to marshall transaction line item", zap.Error(err))
		} else {
			record := &kgo.Record{Topic: "warehouse.stock.transaction.line_item", Value: json}
			f.kafka.Produce(ctx, record, func(_ *kgo.Record, err error) {
				defer wg.Done()
				if err != nil {
					f.logger.For(ctx).Error("Failed to produce event", zap.Error(err))
				}
			})
		}
	}
	wg.Wait()

	return v, err
}

func (f *service) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	v, err := f.repository.GetTransaction(ctx, transactionID)
	return v, err
}
