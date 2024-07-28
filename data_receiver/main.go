package main

import (
	"fmt"
	// "encoding/json"
	"log"
	"net/http"

	"github.com/Fito305/tolling/types"
	// "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gorilla/websocket"
)


func main() {

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	// prod  *kafka.Producer // This is a hard coded dependency
	prod DataProducer // This removed the hard coded dep above.
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p DataProducer
		err error
		kafkaTopic = "obudata"
	)

	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	// Delivery report handler for produced messages
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128), // channel will block after 128.
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
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
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}
}

// In concurrency, a channel will always block when it's full
// The receiver is going to receive the data and it is going to put it
// on a KAFKA queue. We set up Kafka on docker.

// The message channel has a capacity of 128. Once it reaches that it stops receiving data. The channel gets blocked.
