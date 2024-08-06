package main

import (
	"log"

	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/Fito305/tolling/aggregator/client"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

// Transport could be (HTTP, GRPC, kafka) -> attach business logic to this transport
// we are using kafka. But the beautiful thing is that if we wanted to change the transport to 
// for example, JSON, we could just inject it into the interface instead of having to rewrite all
// the code. 

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
