package location

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
	GetLocationsEndpoint   endpoint.Endpoint
	CreateLocationEndpoint endpoint.Endpoint
	DeleteLocationEndpoint endpoint.Endpoint
	GetLocationEndpoint    endpoint.Endpoint
	UpdateLocationEndpoint endpoint.Endpoint
}

//go:embed policies/endpoint.rego
var endpointPolicy string

func NewEndpointSet(s Service, logger log.Factory, tracer opentracing.Tracer) EndpointSet {
	authn := jwt.NewAuthenticator(logger, tracer)
	authz := opa.NewAuthorizor(logger, tracer)

	var getLocationsEndpoint endpoint.Endpoint
	{
		getLocationsEndpoint = makeGetLocationsEndpoint(s, logger, tracer)
		getLocationsEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.location.endpoint.authz.get_locations")(getLocationsEndpoint)
		getLocationsEndpoint = authn.NewMiddleware()(getLocationsEndpoint)
		getLocationsEndpoint = kittracing.TraceServer(tracer, "GetLocationsEndpoint")(getLocationsEndpoint)
	}
	var createLocationEndpoint endpoint.Endpoint
	{
		createLocationEndpoint = makeCreateLocationEndpoint(s, logger, tracer)
		createLocationEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.location.endpoint.authz.create_location")(createLocationEndpoint)
		createLocationEndpoint = authn.NewMiddleware()(createLocationEndpoint)
		createLocationEndpoint = kittracing.TraceServer(tracer, "CreateLocationEndpoint")(createLocationEndpoint)
	}
	var deleteLocationEndpoint endpoint.Endpoint
	{
		deleteLocationEndpoint = makeDeleteLocationEndpoint(s, logger, tracer)
		deleteLocationEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.location.endpoint.authz.delete_location")(deleteLocationEndpoint)
		deleteLocationEndpoint = authn.NewMiddleware()(deleteLocationEndpoint)
		deleteLocationEndpoint = kittracing.TraceServer(tracer, "DeleteLocationEndpoint")(deleteLocationEndpoint)
	}
	var getLocationEndpoint endpoint.Endpoint
	{
		getLocationEndpoint = makeGetLocationEndpoint(s, logger, tracer)
		getLocationEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.location.endpoint.authz.get_location")(getLocationEndpoint)
		getLocationEndpoint = authn.NewMiddleware()(getLocationEndpoint)
		getLocationEndpoint = kittracing.TraceServer(tracer, "GetLocationEndpoint")(getLocationEndpoint)
	}
	var updateLocationEndpoint endpoint.Endpoint
	{
		updateLocationEndpoint = makeUpdateLocationEndpoint(s, logger, tracer)
		updateLocationEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.location.endpoint.authz.update_location")(updateLocationEndpoint)
		updateLocationEndpoint = authn.NewMiddleware()(updateLocationEndpoint)
		updateLocationEndpoint = kittracing.TraceServer(tracer, "UpdateLocationEndpoint")(updateLocationEndpoint)
	}
	return EndpointSet{
		GetLocationsEndpoint:   getLocationsEndpoint,
		CreateLocationEndpoint: createLocationEndpoint,
		DeleteLocationEndpoint: deleteLocationEndpoint,
		GetLocationEndpoint:    getLocationEndpoint,
		UpdateLocationEndpoint: updateLocationEndpoint,
	}
}

// GetLocations

func makeGetLocationsEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Location.GetLocationsEndpoint received request")

		v, err := s.GetLocations(ctx)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// CreateLocation

type CreateLocationEndpointRequestBody CreateLocation

type CreateLocationEndpointRequest struct {
	CreateLocationEndpointRequestBody
}

func makeCreateLocationEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Location.CreateLocationEndpoint received request")

		er := request.(CreateLocationEndpointRequest)
		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Location from endpoint request JSON
		var sr Location
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.CreateLocation(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// DeleteLocation

type DeleteLocationEndpointRequest struct {
	LocationID string
}

func makeDeleteLocationEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Location.DeleteLocationEndpoint received request")
		return nil, nil
	}
}

// GetLocation

type GetLocationEndpointRequest struct {
	LocationID string
}

func makeGetLocationEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Location.GetLocationEndpoint received request")

		er := request.(GetLocationEndpointRequest)
		v, err := s.GetLocation(ctx, er.LocationID)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}

// UpdateLocation

type UpdateLocationEndpointRequestBody UpdateLocation

type UpdateLocationEndpointRequest struct {
	LocationID string
	UpdateLocationEndpointRequestBody
}

func makeUpdateLocationEndpoint(s Service, logger log.Factory, tracer opentracing.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logger.For(ctx).Info("Location.UpdateLocationEndpoint received request")

		er := request.(UpdateLocationEndpointRequest)
		// Convert endpoint request to JSON
		erJSON, err := json.Marshal(er)
		if err != nil {
			return nil, err
		}

		// Create Location from endpoint request JSON
		var sr Location
		json.Unmarshal(erJSON, &sr)

		// Set variables from path parameters
		sr.ID = er.LocationID

		//
		// TODO: Review the code above.
		//       The JSON marshalling isn't ideal.
		//       You should manually construct the struct being passed
		//       to the service from variables in the endpoint request
		//

		v, err := s.UpdateLocation(ctx, &sr)
		if err != nil {
			return &v, err
		}
		return &v, nil

	}
}
