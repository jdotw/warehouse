package item

import (
	"context"
	_ "embed"
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

	getItemHandler := httptransport.NewServer(
		endpoints.GetItemEndpoint,
		decodeGetItemEndpointRequest,
		transport.HTTPEncodeResponse,
		options...,
	)
	r.Handle("/items/{item_id}", getItemHandler).Methods("GET")

}

// GetItem

func decodeGetItemEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var endpointRequest GetItemEndpointRequest

	vars := mux.Vars(r)
	endpointRequest.ItemID = vars["item_id"]

	return endpointRequest, nil
}
