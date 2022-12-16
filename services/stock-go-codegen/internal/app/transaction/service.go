package transaction

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
	"github.com/twmb/franz-go/pkg/kgo"
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
	transaction.logger = f.logger
	transaction.kafka = f.kafka
	v, err := f.repository.CreateTransaction(ctx, transaction)
	return v, err
}

func (f *service) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	v, err := f.repository.GetTransaction(ctx, transactionID)
	return v, err
}
