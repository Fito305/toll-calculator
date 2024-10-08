package aggservice

import (
	"context"
	"fmt"
	"os"

	"github.com/Fito305/tolling/types"
	"github.com/go-kit/kit/log"
)

const basePrice = 3.15 // const go ontop of your file.

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type BasicService struct{
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	// This will print if the wrapped logging does not.
	fmt.Println("this is coming from the internal business logic layer")
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID: obuID,
		TotalDistance: dist,
		TotalAmount: basePrice * dist,
	}
	return inv, nil
}

// NewAggregatorService will construct a complete microservice
// with Logging and instrumentation middleware.
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = newBasicService(NewMemoryStore()) // We are making these onions by wrapping them.
		svc = newLoggingMiddleware(logger)(svc)
		svc = newInstrumentationMiddleware()(svc)
	}
	return svc
}
