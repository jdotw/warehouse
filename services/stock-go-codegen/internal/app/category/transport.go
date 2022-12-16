package category

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

	getCategoriesHandler := httptransport.NewServer(
		endpoints.GetCategoriesEndpoint,
		decodeGetCategoriesEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories", getCategoriesHandler).Methods("GET")

	createCategoryHandler := httptransport.NewServer(
		endpoints.CreateCategoryEndpoint,
		decodeCreateCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories", createCategoryHandler).Methods("POST")

	deleteCategoryHandler := httptransport.NewServer(
		endpoints.DeleteCategoryEndpoint,
		decodeDeleteCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories/{category_id}", deleteCategoryHandler).Methods("DELETE")

	getCategoryHandler := httptransport.NewServer(
		endpoints.GetCategoryEndpoint,
		decodeGetCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories/{category_id}", getCategoryHandler).Methods("GET")

	updateCategoryHandler := httptransport.NewServer(
		endpoints.UpdateCategoryEndpoint,
		decodeUpdateCategoryEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/categories/{category_id}", updateCategoryHandler).Methods("PATCH")

}

// GetCategories

func decodeGetCategoriesEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

// CreateCategory

func decodeCreateCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest CreateCategoryEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	return endpointRequest, nil
}

// DeleteCategory

func decodeDeleteCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest DeleteCategoryEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.CategoryID = vars["category_id"]

	return endpointRequest, nil
}

// GetCategory

func decodeGetCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetCategoryEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.CategoryID = vars["category_id"]

	return endpointRequest, nil
}

// UpdateCategory

func decodeUpdateCategoryEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest UpdateCategoryEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	endpointRequest.CategoryID = vars["category_id"]

	return endpointRequest, nil
}
