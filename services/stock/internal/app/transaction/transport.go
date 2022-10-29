package transaction

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jdotw/go-utils/authn/jwt"
	"github.com/jdotw/go-utils/log"
	"github.com/jdotw/go-utils/transport"
	"github.com/opentracing/opentracing-go"
)

func AddHTTPRoutes(r *mux.Router, endpoints EndpointSet, logger log.Factory, tracer opentracing.Tracer) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(transport.HTTPErrorEncoder),
		httptransport.ServerBefore(jwt.HTTPAuthorizationToContext()),
	}

	getTransactionsHandler := httptransport.NewServer(
		endpoints.GetTransactionsEndpoint,
		decodeGetTransactionsEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/transactions", getTransactionsHandler).Methods("GET")

	createTransactionHandler := httptransport.NewServer(
		endpoints.CreateTransactionEndpoint,
		decodeCreateTransactionEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/transactions", createTransactionHandler).Methods("POST")

	getTransactionHandler := httptransport.NewServer(
		endpoints.GetTransactionEndpoint,
		decodeGetTransactionEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/transactions/{transaction_id}", getTransactionHandler).Methods("GET")

}

// GetTransactions

func decodeGetTransactionsEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetTransactionsEndpointRequest

	return endpointRequest, nil
}

// CreateTransaction

func decodeCreateTransactionEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest CreateTransactionEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	return endpointRequest, nil
}

// GetTransaction

func decodeGetTransactionEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetTransactionEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.TransactionID = vars["transaction_id"]

	return endpointRequest, nil
}
