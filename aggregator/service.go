package main

import (

	"github.com/Fito305/tolling/types"
	"github.com/sirupsen/logrus"
)

const basePrice = 3.15

// The invoicer is going to receive data from the distanceCalculator. It's going to 
// calculate the distance and send it to the invoicer which is going to aggregate
// these distances. This is the invoicer interface.
type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error) // * is a pointer
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuid": distance.OBUID,
		"distance": distance.Value,
		"unix": distance.Unix,
	}).Info("aggregating distance")
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)
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
