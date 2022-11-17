package transaction

import (
	"context"
	_ "embed"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type EndpointSet struct {
	GetTransactionsEndpoint   endpoint.Endpoint
	CreateTransactionEndpoint endpoint.Endpoint
	GetTransactionEndpoint    endpoint.Endpoint
}

//go:embed policies/endpoint.rego
var endpointPolicy string

func NewEndpointSet(s Service, logger log.Factory, tracer opentracing.Tracer) EndpointSet {
	// authn := jwt.NewAuthenticator(logger, tracer)
	// authz := opa.NewAuthorizor(logger, tracer)

	var getTransactionsEndpoint endpoint.Endpoint
	{
		getTransactionsEndpoint = makeGetTransactionsEndpoint(s, logger, tracer)
		// getTransactionsEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.transaction.endpoint.authz.get_transactions")(getTransactionsEndpoint)
		// getTransactionsEndpoint = authn.NewMiddleware()(getTransactionsEndpoint)
		// getTransactionsEndpoint = kittracing.TraceServer(tracer, "GetTransactionsEndpoint")(getTransactionsEndpoint)
	}
	var createTransactionEndpoint endpoint.Endpoint
	{
		createTransactionEndpoint = makeCreateTransactionEndpoint(s, logger, tracer)
		// createTransactionEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.transaction.endpoint.authz.create_transaction")(createTransactionEndpoint)
		// createTransactionEndpoint = authn.NewMiddleware()(createTransactionEndpoint)
		// createTransactionEndpoint = kittracing.TraceServer(tracer, "CreateTransactionEndpoint")(createTransactionEndpoint)
	}
	var getTransactionEndpoint endpoint.Endpoint
	{
		getTransactionEndpoint = makeGetTransactionEndpoint(s, logger, tracer)
		// getTransactionEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.transaction.endpoint.authz.get_transaction")(getTransactionEndpoint)
		// getTransactionEndpoint = authn.NewMiddleware()(getTransactionEndpoint)
		// getTransactionEndpoint = kittracing.TraceServer(tracer, "GetTransactionEndpoint")(getTransactionEndpoint)
	}
	return EndpointSet{
		GetTransactionsEndpoint:   getTransactionsEndpoint,
		CreateTransactionEndpoint: createTransactionEndpoint,
		GetTransactionEndpoint:    getTransactionEndpoint,
	}
}

// GetTransactions

func makeGetTransactionsEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Transaction.GetTransactionsEndpoint received request")

		v, err := s.GetTransactions(ctx)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// CreateTransaction

type CreateTransactionEndpointRequestBody CreateTransaction

type CreateTransactionEndpointRequest struct {
	CreateTransactionEndpointRequestBody
}

func makeCreateTransactionEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		er := request.(CreateTransactionEndpointRequest)

		// Create the Transaction
		var sr Transaction
		sr.LocationID = er.LocationID
		if er.Timestamp != nil {
			sr.Timestamp = *er.Timestamp
		} else {
			sr.Timestamp = time.Now()
		}
		sr.Items = make([]TransactionLineItem, len(er.Items))
		for i, el := range er.Items {
			var sl TransactionLineItem
			sl.ItemID = el.ItemID
			sl.Quantity = el.Quantity
			sr.Items[i] = sl
		}

		v, err := s.CreateTransaction(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil
	}
}

// GetTransaction

type GetTransactionEndpointRequest struct {
	TransactionID string
}

func makeGetTransactionEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Transaction.GetTransactionEndpoint received request")

		er := request.(GetTransactionEndpointRequest)
		v, err := s.GetTransaction(ctx, er.TransactionID)
		if err != nil {
			logger.For(ctx).Error("GetTransactionEndpoint failed", zap.Error(err))
			return &v, err
		}
		return &v, nil

	}
}
