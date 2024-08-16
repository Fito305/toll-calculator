package aggendpoint

import (
	"context"

	"github.com/Fito305/tolling/types"
	"github.com/Fito305/tolling/go-kit-example/aggsvc/aggservice"
	"github.com/go-kit/kit/endpoint"
)

// Used to implement multiple endpoints otherwise you would need 10 arguments.
type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

// The AggregateRequest on the other hand, is what the transport
// needs to serialize it based on the transport mechanism. Then we need to message that
// into our business logic type Distance struct in types.go.
// you can embed the Distance struct here (less redable) but it uses the same values.
type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{OBUID: obuID})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID: result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount: result.TotalAmount,
	}, nil
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{ // _ is the resp (response)
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	return err
}

func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix: req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeConcatEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{ // We ask our business logic to calculate this invoice (inv) and then we are going to return a response in this struct. And that's going to be serialized JSON.
			Err: err,
			OBUID: inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount: inv.TotalAmount,
		}, nil
	}
}
