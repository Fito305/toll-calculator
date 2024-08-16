// The aggregator directory is the invoicer
package main

import (
	// "context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	// "time"
	// "strconv"

	"github.com/Fito305/tolling/types"
	"github.com/joho/godotenv"
	// "github.com/Fito305/tolling/aggregator/client"
	"google.golang.org/grpc"
	// "github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// to use the .env file enviroment variables.
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store = makeStore()
		// store = NewMemoryStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc)) // *grpcListenAddr dereference. You need to put it in a `go` routine. But need to have a mechanic to close it.
	}()
	// if you see a star like this it means *listenAddr - we are dereferencing it with *.
	log.Fatal(makeHTTPTransport(httpListenAddr, svc)) // parameters must be passed in the same order as the func definition.
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
	var (
		aggMetricHandler = newHTTPMetricsHandler("aggregate")
		invMetricHandler = newHTTPMetricsHandler("invoice")
		aggregateHandler = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
		invoiceHandler   = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	)
	http.HandleFunc("/invoice", invoiceHandler)
	http.HandleFunc("/aggregate", aggregateHandler)
	// http.HandleFunc("/invoice", invMetricHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port ", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
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
