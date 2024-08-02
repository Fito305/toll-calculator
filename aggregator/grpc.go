package main

import (
	"context"

	"github.com/Fito305/tolling/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
		}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuID), // conversions happen here.
		Value: req.Value,	// via the req parameter.
		Unix: req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)	// (1)
	// Why can we just return it? Because AggregateDistance is returning an error.
}



// We are going to have an HTTP server and each time you hit request with the JSON it is
// going to use the service and aggregate the distance. But now we are going to have a GRPC 
// transport so it is not going to work with our HTTP server so we need to make a 
// GRPC server. AggregateDistance() is a method of the server.

// svc stands for service.

// So we are going to create a new GRPC server, we are going to make a constructor which 
// is called NewGRPCServer() so we can put in our Aggregator. So we can basically say here (1), 
// s.svc.AggregateDistance().

// We can't use the same types because there is a big difference between the ``transport layer`` 
// and ``business layer`` (very important). The ``business layer`` is going to have it's business 
// layer type. It's the main type that everyone needs to convert to. No matter what transport
// it is, everybody will convert their request and response type to the business layer type.
// The business layer is central it decides.
// The ``transport layer``, even if it's a JSON, we can just use the types.Distance because it is the same.
// For the JSON it is already types.Distance so no conversion needed.
// !!! But for the GRPC we use the types.AggregateRequest in to ptypes.proto file. Because it is a request
// and it needs to convert it to a types.Distance.

// An maybe you'll have for example a ```Webpack``` transport then it'll go to type.WEBpack then convert 
// to types.Distance. So everybody has to convert to a types.Distance.

// Normally a rpc request means that it is a request and a response. "Hey give me the username" and it 
// returns the username. But in our case it is AggregateDistance, which we send but the 
// only thing we care about is if there is an error or not. We don't have any result and that is why
// we use the type.None in the return of Aggregate()
