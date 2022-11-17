package transaction

import (
	"context"
	_ "embed"
)

type Repository interface {
	GetTransactions(ctx context.Context) (*[]Transaction, error)
	CreateTransaction(ctx context.Context, transaction *Transaction) (*Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*Transaction, error)
}
