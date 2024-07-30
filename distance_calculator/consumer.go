package main

import (
	"encoding/json"
	"time"
	
	"github.com/Fito305/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/Fito305/tolling/aggregator/client"
	"github.com/sirupsen/logrus"
)

// This can also be called KafkaTransport.
// Because the business logic will be abstract.
// We get to wrap things in different layers that are independant from each other.
// This is not being taught.
type KafkaConsumer struct {
	consumer *kafka.Consumer
	isRunning bool
	calcService CalculatorServicer
	aggClient *client.Client // The start * is a pointer // We make the client so all the services can import this client.
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)


	return &KafkaConsumer{
		consumer: c,
		calcService: svc,
		aggClient: aggClient, 
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {

	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue // if we log out, and break the loop we won't recieve anything. 
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err) // ***!!! Very important, in development, when we first make these things we use JSOn because it is easy to debug. But once we go to production, once the whole architecture is set up and all tests are working, **Then we swap out to **ProtoBuffers to make it faster. Once everything is implemented. 
			continue // This is important we are going ot continue like above.
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}
		req := types.Distance{
			Value: distance,
			Unix: time.Now().UnixNano(),
			OBUID: data.OBUID,
		}
		if err := c.aggClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error:", err)
			continue // Does not matter unless you put some more logic below.
		}
	}
}
