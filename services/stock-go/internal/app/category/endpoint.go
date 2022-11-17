package category

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
)

type EndpointSet struct {
	GetCategoriesEndpoint  endpoint.Endpoint
	CreateCategoryEndpoint endpoint.Endpoint
	DeleteCategoryEndpoint endpoint.Endpoint
	GetCategoryEndpoint    endpoint.Endpoint
	UpdateCategoryEndpoint endpoint.Endpoint
}

//go:embed policies/endpoint.rego
var endpointPolicy string

func NewEndpointSet(s Service, logger log.Factory, tracer opentracing.Tracer) EndpointSet {
	authn := jwt.NewAuthenticator(logger, tracer)
	authz := opa.NewAuthorizor(logger, tracer)

	var getCategoriesEndpoint endpoint.Endpoint
	{
		getCategoriesEndpoint = makeGetCategoriesEndpoint(s, logger, tracer)
		getCategoriesEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.category.endpoint.authz.get_categories")(getCategoriesEndpoint)
		getCategoriesEndpoint = authn.NewMiddleware()(getCategoriesEndpoint)
		getCategoriesEndpoint = kittracing.TraceServer(tracer, "GetCategoriesEndpoint")(getCategoriesEndpoint)
	}
	var createCategoryEndpoint endpoint.Endpoint
	{
		createCategoryEndpoint = makeCreateCategoryEndpoint(s, logger, tracer)
		createCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.category.endpoint.authz.create_category")(createCategoryEndpoint)
		createCategoryEndpoint = authn.NewMiddleware()(createCategoryEndpoint)
		createCategoryEndpoint = kittracing.TraceServer(tracer, "CreateCategoryEndpoint")(createCategoryEndpoint)
	}
	var deleteCategoryEndpoint endpoint.Endpoint
	{
		deleteCategoryEndpoint = makeDeleteCategoryEndpoint(s, logger, tracer)
		deleteCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.category.endpoint.authz.delete_category")(deleteCategoryEndpoint)
		deleteCategoryEndpoint = authn.NewMiddleware()(deleteCategoryEndpoint)
		deleteCategoryEndpoint = kittracing.TraceServer(tracer, "DeleteCategoryEndpoint")(deleteCategoryEndpoint)
	}
	var getCategoryEndpoint endpoint.Endpoint
	{
		getCategoryEndpoint = makeGetCategoryEndpoint(s, logger, tracer)
		getCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.category.endpoint.authz.get_category")(getCategoryEndpoint)
		getCategoryEndpoint = authn.NewMiddleware()(getCategoryEndpoint)
		getCategoryEndpoint = kittracing.TraceServer(tracer, "GetCategoryEndpoint")(getCategoryEndpoint)
	}
	var updateCategoryEndpoint endpoint.Endpoint
	{
		updateCategoryEndpoint = makeUpdateCategoryEndpoint(s, logger, tracer)
		updateCategoryEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.category.endpoint.authz.update_category")(updateCategoryEndpoint)
		updateCategoryEndpoint = authn.NewMiddleware()(updateCategoryEndpoint)
		updateCategoryEndpoint = kittracing.TraceServer(tracer, "UpdateCategoryEndpoint")(updateCategoryEndpoint)
	}
	return EndpointSet{
		GetCategoriesEndpoint:  getCategoriesEndpoint,
		CreateCategoryEndpoint: createCategoryEndpoint,
		DeleteCategoryEndpoint: deleteCategoryEndpoint,
		GetCategoryEndpoint:    getCategoryEndpoint,
		UpdateCategoryEndpoint: updateCategoryEndpoint,
	}
}

// GetCategories


func makeGetCategoriesEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Category.GetCategoriesEndpoint received request")

		v, err := s.GetCategories(ctx)
		if err != nil {
			return &v, err
		}
		return &v, nil
	}
}

// CreateCategory

type CreateCategoryEndpointRequestBody CreateCategory

type CreateCategoryEndpointRequest struct {
	CreateCategoryEndpointRequestBody
}

func makeCreateCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Category.CreateCategoryEndpoint received request")

		er := request.(CreateCategoryEndpointRequest)
		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Category from endpoint request JSON
		var sr Category
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.CreateCategory(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// DeleteCategory

type DeleteCategoryEndpointRequest struct {
	CategoryID string
}

func makeDeleteCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Category.DeleteCategoryEndpoint received request")

		er := request.(DeleteCategoryEndpointRequest)
		err := s.DeleteCategory(ctx, er.CategoryID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

// GetCategory

type GetCategoryEndpointRequest struct {
	CategoryID string
}

func makeGetCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Category.GetCategoryEndpoint received request")

		er := request.(GetCategoryEndpointRequest)
		v, err := s.GetCategory(ctx, er.CategoryID)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// UpdateCategory

type UpdateCategoryEndpointRequestBody UpdateCategory

type UpdateCategoryEndpointRequest struct {
	CategoryID string
	UpdateCategoryEndpointRequestBody
}

func makeUpdateCategoryEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Category.UpdateCategoryEndpoint received request")

		er := request.(UpdateCategoryEndpointRequest)
		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Category from endpoint request JSON
		var sr Category
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters
		sr.ID = er.CategoryID

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.UpdateCategory(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}
