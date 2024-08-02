// The aggregator directory is the invoicer
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"log"
	"time"
	"strconv"

	"github.com/Fito305/tolling/types"
	"github.com/Fito305/tolling/aggregator/client"
	"google.golang.org/grpc"
	// "github.com/sirupsen/logrus"
)

func main() {
	// You can pass in a `flag` in the command line to change the address `--listenaddr <port>`
	httpListenAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP transport server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen address of the GRPC transport server")
	flag.Parse() // you have to parse it.
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go func () {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc)) // *grpcListenAddr dereference. You need to put it in a `go` routine. But need to have a mechanic to close it.
	}()
	time.Sleep(time.Second * 5)
	c, err := client.NewGRPCClient(*grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix: time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	// if you see a star like this it means *listenAddr - we are dereferencing it with *.
	log.Fatal(makeHTTPTransport(*httpListenAddr, svc)) // parameters must be passed in the same order as the func definition.
	// Make a transporter
}

// GRPC server. Transport.
func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port ", listenAddr)
	// We need to make a TCP listener first.
	ln, err := net.Listen("tcp", listenAddr) // "tcp" has to be in lowercase.
	if err != nil {
		return err
	}
	defer ln.Close() // close the 'go' routine in main() 
	// Make a new GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...) // ... is an elipses
	// Register (OUR) GRPC server implementation to the GRPC package
	// This serves the GRPC request.
	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc)) // server is the native server & NEWGRPCServer is the GRPC server.
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc)) // You cannot attach the http transport to the Aggregator interface that is something you cannot do. You can do a http.HandleFunc
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddr, nil)
}

// This is the HTTP JSON server. Transport.
func handleGetInvoice(svc Aggregator) http.HandlerFunc { // The decorator allows us to intergrate interface Aggregator for http use cases.
	return func(w http.ResponseWriter, r *http.Request) { // !!! DECORATOR PATTERN !!!
	// fmt.Println(r.URL.Query()) // how to infer the url to get it logged. Query() gives you a map and URL alone jsut the URL string.
	values, ok := r.URL.Query()["obu"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
		return
	}
	obuID, err := strconv.Atoi(values[0])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
		return
	}
	invoice, err := svc.CalculateInvoice(obuID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}) // we do err.Error() because if you want a string of error you have to call Error()
		return 
	}
	writeJSON(w, http.StatusOK, invoice)
  }
}

// A transport
func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // !Decorator Pattern!
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v) // Why v? it's used for interfaces when it's an unknown varible of type any.
}

// This invoicer is going to have a transport.
// How are we going to reach this invoicer?
// Are we going to use JSON or ProtoBuffers?
// We are going to use both!
// Why? You need to set up your Micro Service `transport independant`.
// We first start off with JSON it's easier to debug. And that is how companies start.
// And once everything is set up we can add another transport that is going to
// be the Protobuffers. Knowing how to implement Protobuffer is very useful.

// Very important: if you are making Mircro Services, in a company. You are not going to just make your service and call
// it a day. You are also going to make a client for your Micro Service. A client which is a simple package that
// people can import to interact with your Micro Service.


// makeGRPCTransport() in this function, we are going to make a TCP listener. 
// Then we make the server for GRPC. The server from the package itself. From the GRPC 
// package. And then we need to register our Aggregator server.
