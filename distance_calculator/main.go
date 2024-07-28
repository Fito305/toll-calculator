package main

import (
	"log"

	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const kafkaTopic = "obudata"

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

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
