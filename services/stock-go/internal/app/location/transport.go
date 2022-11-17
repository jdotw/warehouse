package location

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

	getLocationsHandler := httptransport.NewServer(
		endpoints.GetLocationsEndpoint,
		decodeGetLocationsEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/locations", getLocationsHandler).Methods("GET")

	createLocationHandler := httptransport.NewServer(
		endpoints.CreateLocationEndpoint,
		decodeCreateLocationEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/locations", createLocationHandler).Methods("POST")

	deleteLocationHandler := httptransport.NewServer(
		endpoints.DeleteLocationEndpoint,
		decodeDeleteLocationEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/locations/{location_id}", deleteLocationHandler).Methods("DELETE")

	getLocationHandler := httptransport.NewServer(
		endpoints.GetLocationEndpoint,
		decodeGetLocationEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/locations/{location_id}", getLocationHandler).Methods("GET")

	updateLocationHandler := httptransport.NewServer(
		endpoints.UpdateLocationEndpoint,
		decodeUpdateLocationEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/locations/{location_id}", updateLocationHandler).Methods("PATCH")

}

// GetLocations

func decodeGetLocationsEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

// CreateLocation

func decodeCreateLocationEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest CreateLocationEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	return endpointRequest, nil
}

// DeleteLocation

func decodeDeleteLocationEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest DeleteLocationEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.LocationID = vars["location_id"]

	return endpointRequest, nil
}

// GetLocation

func decodeGetLocationEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetLocationEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.LocationID = vars["location_id"]

	return endpointRequest, nil
}

// UpdateLocation

func decodeUpdateLocationEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest UpdateLocationEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&endpointRequest); err != nil {
		return nil, err
	}
	vars := mux.Vars(r)
	endpointRequest.LocationID = vars["location_id"]

	return endpointRequest, nil
}
