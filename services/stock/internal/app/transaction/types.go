// DO NOT EDIT
// This file was code-generated by  version
// It is expected that this file will be re-generated and overwitten to
// adapt to changes in the OpenAPI spec that was used to generate it

package transaction

import (
	_ "embed"
	"time"
)

// Item in a new Transaction
type CreateTransactionItem struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// Create Transaction
type CreateTransactionRequest struct {
	Items     []CreateTransactionRequestItem `json:"items"`
	Timestamp *time.Time                     `json:"timestamp,omitempty"`
}

// CreateTransactionRequestItem defines model for []create_transaction_request_item.
type CreateTransactionRequestItem CreateTransactionItem

// HTTPError defines model for HTTPError.
type HTTPError struct {
	Message *string `json:"message,omitempty"`
}

// Transaction
type Transaction struct {
	ID        string            `gorm:"primaryKey;unique;type:uuid;default:uuid_generate_v4();" json:"id"`
	Items     []TransactionItem `json:"items"`
	Timestamp time.Time         `gorm:"not null" json:"timestamp"`
}

// TransactionItem defines model for []transaction_item.
type TransactionItem TransactionLineItem

// Item in a Transaction
type TransactionLineItem struct {
	ID            string `gorm:"primaryKey;unique;type:uuid;default:uuid_generate_v4();" json:"id"`
	ItemID        string `gorm:"not null" json:"item_id"`
	Quantity      int    `json:"quantity"`
	TransactionID string `gorm:"not null" json:"transaction_id"`
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
