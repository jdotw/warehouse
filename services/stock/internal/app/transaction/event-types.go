// DO NOT EDIT
// This file was code-generated by  version
// It is expected that this file will be re-generated and overwitten to
// adapt to changes in the OpenAPI spec that was used to generate it

package transaction

import (
	_ "embed"
)

// Transaction Line Item Created Event
type TransactionLineItemCreatedEvent struct {
	LocationID string `json:"location_id"`
	LineItem TransactionLineItem `json:"line_item"`
}

