package transaction

import (
	_ "embed"
)

// Transaction Line Item Created Event
type TransactionLineItemCreatedEvent struct {
	LineItem TransactionLineItem `json:"line_item"`
}
