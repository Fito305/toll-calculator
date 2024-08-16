package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Fito305/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)
// we can't catch a error because HandlerFunc does not return an error. You as the programmer has to implement error handling. Its added via the HTTPFunc type.
// This and APIError and the Error func below inplements it.
type HTTPFunc func(http.ResponseWriter, *http.Request) error // tp be able to use error handling for Handler Func.

type APIError struct {
	Code int
	Err  error
}

// Error implements the error interface.
func (e APIError) Error() string {
	return e.Err.Error()
}

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}
		}
	}
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
		Name:      "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name:      "latency",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &HTTPMetricHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
		errCounter: errCounter,
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	// gets call once
	return func(w http.ResponseWriter, r *http.Request) error { // decorator pattern
		var err error
		// gets called for each request. That's why you put the counter in here.
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err": err,
			}).Info() // have to put debug() or info() at the end of logrus or it won't work.
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())
		err = next(w, r)
		return err
	}

}

// This is the HTTP JSON server. Transport.
func handleGetInvoice(svc Aggregator) HTTPFunc { // The decorator allows us to intergrate interface Aggregator for http use cases.
	return func(w http.ResponseWriter, r *http.Request) error { // !!! DECORATOR PATTERN !!!
		// fmt.Println(r.URL.Query()) // how to infer the url to get it logged. Query() gives you a map and URL alone jsut the URL string.
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing OBU id"),
			}
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU id %s", values[0]),
			}
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, invoice)
	}
}

// A transport
func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error { // !Decorator Pattern!
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("method not supported (%s)", r.Method),
			}
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Errorf("failed to decode the response body %s", err),
			}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}
