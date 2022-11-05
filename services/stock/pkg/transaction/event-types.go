package transaction

import (
	_ "embed"
)

// Transaction Line Item Created Event
type TransactionLineItemCreatedEvent struct {
	LocationID string              `json:"location_id"`
	LineItem   TransactionLineItem `json:"line_item"`
}
