package item

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	kittracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/jdotw/go-utils/authn/jwt"
	"github.com/jdotw/go-utils/authz/opa"
	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type EndpointSet struct {
	GetItemsInCategoryEndpoint   endpoint.Endpoint
	CreateItemInCategoryEndpoint endpoint.Endpoint
	DeleteItemEndpoint           endpoint.Endpoint
	GetItemEndpoint              endpoint.Endpoint
	UpdateItemEndpoint           endpoint.Endpoint
}

//go:embed policies/endpoint.rego
var endpointPolicy string

func NewEndpointSet(s Service, logger log.Factory, tracer opentracing.Tracer) EndpointSet {
	authn := jwt.NewAuthenticator(logger, tracer)
	authz := opa.NewAuthorizor(logger, tracer)

	var getItemsInCategoryEndpoint endpoint.Endpoint
	{
		getItemsInCategoryEndpoint = makeGetItemsInCategoryEndpoint(s, logger, tracer)
		getItemsInCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.get_items_in_category")(getItemsInCategoryEndpoint)
		getItemsInCategoryEndpoint = authn.NewMiddleware()(getItemsInCategoryEndpoint)
		getItemsInCategoryEndpoint = kittracing.TraceServer(tracer, "GetItemsInCategoryEndpoint")(getItemsInCategoryEndpoint)
	}
	var createItemInCategoryEndpoint endpoint.Endpoint
	{
		createItemInCategoryEndpoint = makeCreateItemInCategoryEndpoint(s, logger, tracer)
		createItemInCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.create_item_in_category")(createItemInCategoryEndpoint)
		createItemInCategoryEndpoint = authn.NewMiddleware()(createItemInCategoryEndpoint)
		createItemInCategoryEndpoint = kittracing.TraceServer(tracer, "CreateItemInCategoryEndpoint")(createItemInCategoryEndpoint)
	}
	var deleteItemEndpoint endpoint.Endpoint
	{
		deleteItemEndpoint = makeDeleteItemEndpoint(s, logger, tracer)
		deleteItemEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.delete_item")(deleteItemEndpoint)
		deleteItemEndpoint = authn.NewMiddleware()(deleteItemEndpoint)
		deleteItemEndpoint = kittracing.TraceServer(tracer, "DeleteItemEndpoint")(deleteItemEndpoint)
	}
	var getItemEndpoint endpoint.Endpoint
	{
		getItemEndpoint = makeGetItemEndpoint(s, logger, tracer)
		getItemEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.get_item")(getItemEndpoint)
		getItemEndpoint = authn.NewMiddleware()(getItemEndpoint)
		getItemEndpoint = kittracing.TraceServer(tracer, "GetItemEndpoint")(getItemEndpoint)
	}
	var updateItemEndpoint endpoint.Endpoint
	{
		updateItemEndpoint = makeUpdateItemEndpoint(s, logger, tracer)
		updateItemEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.update_item")(updateItemEndpoint)
		updateItemEndpoint = authn.NewMiddleware()(updateItemEndpoint)
		updateItemEndpoint = kittracing.TraceServer(tracer, "UpdateItemEndpoint")(updateItemEndpoint)
	}
	return EndpointSet{
		GetItemsInCategoryEndpoint:   getItemsInCategoryEndpoint,
		CreateItemInCategoryEndpoint: createItemInCategoryEndpoint,
		DeleteItemEndpoint:           deleteItemEndpoint,
		GetItemEndpoint:              getItemEndpoint,
		UpdateItemEndpoint:           updateItemEndpoint,
	}
}

// GetItemsInCategory

type GetItemsInCategoryEndpointRequest struct {
	CategoryID string
}

func makeGetItemsInCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		er := request.(GetItemsInCategoryEndpointRequest)
		v, err := s.GetItemsInCategory(ctx, er.CategoryID)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// CreateItemInCategory

type CreateItemInCategoryEndpointRequestBody CreateItemInCategory

type CreateItemInCategoryEndpointRequest struct {
	CategoryID string
	CreateItemInCategoryEndpointRequestBody
}

func makeCreateItemInCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		er := request.(CreateItemInCategoryEndpointRequest)
		logger.For(ctx).Info("Item.CreateItemInCategoryEndpoint received request", zap.String("category_id", er.CategoryID), zap.String("name", er.Name), zap.Int("upc", er.UPC))

		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Item from endpoint request JSON
		var sr Item
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters
		sr.CategoryID = er.CategoryID

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.CreateItemInCategory(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// DeleteItem

type DeleteItemEndpointRequest struct {
	ItemID string
}

func makeDeleteItemEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Item.DeleteItemEndpoint received request")
		er := request.(DeleteItemEndpointRequest)
		err := s.DeleteItem(ctx, er.ItemID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// GetItem

type GetItemEndpointRequest struct {
	ItemID string
}

func makeGetItemEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Item.GetItemEndpoint received request")

		er := request.(GetItemEndpointRequest)
		v, err := s.GetItem(ctx, er.ItemID)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// UpdateItem

type UpdateItemEndpointRequestBody UpdateItem

type UpdateItemEndpointRequest struct {
	ItemID string
	UpdateItemEndpointRequestBody
}

func makeUpdateItemEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Item.UpdateItemEndpoint received request")

		er := request.(UpdateItemEndpointRequest)
		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Item from endpoint request JSON
		var sr Item
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters
		sr.ID = er.ItemID

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.UpdateItem(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}
