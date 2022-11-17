package item

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

	getItemsInCategoryHandler := httptransport.NewServer(
		endpoints.GetItemsInCategoryEndpoint,
		decodeGetItemsInCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories/{category_id}/items", getItemsInCategoryHandler).Methods("GET")

	createItemInCategoryHandler := httptransport.NewServer(
		endpoints.CreateItemInCategoryEndpoint,
		decodeCreateItemInCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories/{category_id}/items", createItemInCategoryHandler).Methods("POST")

	deleteItemHandler := httptransport.NewServer(
		endpoints.DeleteItemEndpoint,
		decodeDeleteItemEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/items/{item_id}", deleteItemHandler).Methods("DELETE")

	getItemHandler := httptransport.NewServer(
		endpoints.GetItemEndpoint,
		decodeGetItemEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/items/{item_id}", getItemHandler).Methods("GET")

	updateItemHandler := httptransport.NewServer(
		endpoints.UpdateItemEndpoint,
		decodeUpdateItemEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/items/{item_id}", updateItemHandler).Methods("PATCH")

}

// GetItemsInCategory

func decodeGetItemsInCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetItemsInCategoryEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.CategoryID = vars["category_id"]

	return endpointRequest, nil
}

// CreateItemInCategory

func decodeCreateItemInCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest CreateItemInCategoryEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	endpointRequest.CategoryID = vars["category_id"]

	return endpointRequest, nil
}

// DeleteItem

func decodeDeleteItemEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest DeleteItemEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.ItemID = vars["item_id"]

	return endpointRequest, nil
}

// GetItem

func decodeGetItemEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetItemEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.ItemID = vars["item_id"]

	return endpointRequest, nil
}

// UpdateItem

func decodeUpdateItemEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest UpdateItemEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	endpointRequest.ItemID = vars["item_id"]

	return endpointRequest, nil
}
