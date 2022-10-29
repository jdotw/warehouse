package transaction

import (
	"context"
	_ "embed"

	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type Service interface {
	GetTransactions(ctx context.Context) (*[]Transaction, error)
	CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository, logger log.Factory, tracer opentracing.Tracer) Service {
	var svc Service
	{
		svc = &service{
			repository: repository,
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
	return v, err
}

func (f *service) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	v, err := f.repository.GetTransaction(ctx, transactionID)
	return v, err
}
