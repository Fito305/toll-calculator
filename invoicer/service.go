package main

import (

	
)

// The invoicer is going to receive data from the distanceCalculator. It's going to 
// calculate the distance and send it to the invoicer which is going to aggregate
// these distances.
type Invoicer interface {
	AggregateDistance()

}
