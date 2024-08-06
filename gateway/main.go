package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"
	"net/http"

	"github.com/Fito305/tolling/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error // ! DECORATOR PATTERN/ HIHG ORDER FUNCTIONS !

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listne address of the HTTP address.") // Flag makes it so you can parse the flag from the command line.
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http.//localhost:3000", "the listen address of the aggregator service")
	flag.Parse()
	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr) // Endpoint o fthe aggregator service.
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice)) // wrap it to return an error.
	logrus.Infof("gateway HTTP running on port: %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil)) // *listenAddr is being dereference. Because the variable listenAddr is a pointer above so we need to deregference it.
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error { // same as apiFunc above that's how you get this to return an error.
	// Access to the agg client
	inv, err := h.client.GetInvoice(context.Background(), 4984)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

// Helper function
func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// Converter that allows handleGetInvoice to return an error.
func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri": r.RequestURI,
			}).Info("REQ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

// Start simple, test, add more, test, and keep building and improving your own framework.
