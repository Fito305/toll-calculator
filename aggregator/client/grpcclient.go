package client

import (

	"github.com/Fito305/tolling/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient // It's embedded so we have direct access to everything the AggregatorClient has.
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		AggregatorClient: c, // We are embedding AggregatorClient here.
	}, nil
}
