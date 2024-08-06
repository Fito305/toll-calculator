package client

import (
	"context"


	"github.com/Fito305/tolling/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	client types.AggregatorClient 
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client: c, // We are embedding AggregatorClient here.
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}
