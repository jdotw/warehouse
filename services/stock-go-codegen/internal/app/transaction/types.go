// DO NOT EDIT
// This file was code-generated by  version
// It is expected that this file will be re-generated and overwitten to
// adapt to changes in the OpenAPI spec that was used to generate it

package transaction

import (
	"context"
	_ "embed"
	"time"

	"github.com/jdotw/go-utils/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Item in a new Transaction
type CreateTransactionItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// Create Transaction
type CreateTransactionRequest struct {
	Items      []CreateTransactionRequestItem `json:"items"`
	LocationID string                         `json:"location_id"`
	Timestamp  *time.Time                     `json:"timestamp,omitempty"`
}

// CreateTransactionRequestItem defines model for []create_transaction_request_item.
type CreateTransactionRequestItem CreateTransactionItem

// HTTPError defines model for HTTPError.
type HTTPError struct {
	Message *string `json:"message,omitempty"`
}

// Transaction
type Transaction struct {
	ID         string            `gorm:"primaryKey;unique;type:uuid;default:uuid_generate_v4();" json:"id"`
	Items      []TransactionLineItem `json:"items"`
	LocationID string            `gorm:"not null" json:"location_id"`
	Timestamp  time.Time         `gorm:"not null" json:"timestamp"`

	// Internal
	kafka      *kgo.Client
	logger     log.Factory
	ctx context.Context
}

// Item in a Transaction
type TransactionLineItem struct {
	ID             string `gorm:"primaryKey;unique;type:uuid;default:uuid_generate_v4();" json:"id"`
	ItemID         string `gorm:"not null" json:"item_id"`
	Quantity       int    `json:"quantity"`
	SequenceNumber int    `gorm:"unique;type:uint;autoIncrement;" json:"sequence_number"`
	TransactionID  string `gorm:"not null" json:"transaction_id"`
}

// BadRequestError defines model for BadRequestError.
type BadRequestError HTTPError

// ForbiddenError defines model for ForbiddenError.
type ForbiddenError HTTPError

// InternalServerError defines model for InternalServerError.
type InternalServerError HTTPError

// NotFoundError defines model for NotFoundError.
type NotFoundError HTTPError

// UnauthorizedError defines model for UnauthorizedError.
type UnauthorizedError HTTPError

// Create Transaction
type CreateTransaction CreateTransactionRequest
