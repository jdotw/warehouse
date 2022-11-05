// DO NOT EDIT
// This file was code-generated by  version
// It is expected that this file will be re-generated and overwitten to
// adapt to changes in the OpenAPI spec that was used to generate it

package item

import (
	_ "embed"
)

// HTTPError defines model for HTTPError.
type HTTPError struct {
	Message *string `json:"message,omitempty"`
}

// Stock on Hand for an Item
type ItemStockOnHand struct {
	ID          string `gorm:"primaryKey;unique;type:uuid;" json:"id"`
	StockOnHand int    `gorm:"not null" json:"stock-on-hand"`
}

// BadRequestError defines model for BadRequestError.
type BadRequestError HTTPError

// InternalServerError defines model for InternalServerError.
type InternalServerError HTTPError

// NotFoundError defines model for NotFoundError.
type NotFoundError HTTPError
