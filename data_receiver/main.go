package main

import (
	"fmt"
    "json"
	"log"
	"net/http"

	"github.com/Fito305/tolling/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
)

const kafkaTopic = "obudata"

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := kafkaTopic
	for i := 0; i < 10; i++ {
	}

	recv, err := NewDataReceiver()

	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  *kafka.Producer
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128), // channel will block after 128.
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
    b, err := json.Marshal(data)
    if err != nil {
        return err
    }
	err = dr.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b, //[]byte("testing producing"), replaced by b
	}, nil)
    return err
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected - client connected!")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue // If one truck sends incorrect obudata you cannot close because no other trucks will be able to send data.
		}
		fmt.Printf("received OBU data from [%d] :: <lat %.2f, long %.2f> \n", data.OBUID, data.Lat, data.Long)
		// dr.msgch <- data // we are piping in the data here. Over and over again.
		// The message channel has a capacity of 128. Once it reaches that it stops receiving data. The channel gets blocked.
	}
}

// In concurrency, a channel will always block when it's full
// The receiver is going to receive the data and it is going to put it
// on a KAFKA queue. We set up Kafka on docker.
