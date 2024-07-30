// The aggregator directory is the invoicer
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Fito305/tolling/types"
	// "github.com/sirupsen/logrus"
)

func main() {
	// You can pass in a `flag` in the command line to change the address `--listenaddr <port>`
	listenAddr := flag.String("listenAddr", ":3000", "the listen address of the HTTP transport server")
	flag.Parse() // you have to parse it.
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	// if you see a star like this it means *listenAddr - we are dereferencing it with *.
	makeHTTPTransport(*listenAddr, svc) // parameters must be passed in the same order as the func definition.
	// Make a transporter
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc)) // You cannot attach the http transport to the Aggregator interface that is something you cannot do. You can do a http.HandleFunc
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

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
