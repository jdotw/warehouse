package item

import (
	"context"
	_ "embed"

	"github.com/go-kit/kit/endpoint"
	kittracing "github.com/go-kit/kit/tracing/opentracing"
	"github.com/jdotw/go-utils/authn/jwt"
	"github.com/jdotw/go-utils/authz/opa"
	"github.com/jdotw/go-utils/log"
	"github.com/opentracing/opentracing-go"
)

type EndpointSet struct {
	GetItemEndpoint endpoint.Endpoint
}

//go:embed policies/endpoint.rego
var endpointPolicy string

func NewEndpointSet(s Service, logger log.Factory, tracer opentracing.Tracer) EndpointSet {
	authn := jwt.NewAuthenticator(logger, tracer)
	authz := opa.NewAuthorizor(logger, tracer)

	var getItemEndpoint endpoint.Endpoint
	{
		getItemEndpoint = makeGetItemEndpoint(s, logger, tracer)
		getItemEndpoint = authz.NewInProcessMiddleware(endpointPolicy, "data.item.endpoint.authz.get_item")(getItemEndpoint)
		getItemEndpoint = authn.NewMiddleware()(getItemEndpoint)
		getItemEndpoint = kittracing.TraceServer(tracer, "GetItemEndpoint")(getItemEndpoint)
	}
	return EndpointSet{
		GetItemEndpoint: getItemEndpoint,
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
