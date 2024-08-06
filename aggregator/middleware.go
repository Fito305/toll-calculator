package main

import (
	"time"

	"github.com/Fito305/tolling/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next Aggregator
}
 
func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{     // Promauto registers all your metrics automatically.
		Namespace: "aggregator_request_counter",
		Name: "aggregator",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{     
		Namespace: "aggregator_request_counter",
		Name: "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name: "aggregator",
		Buckets: []float64{0.1, 0.5, 1}, // The buckets is for the histogram.
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name: "calculate",
		Buckets: []float64{0.1, 0.5, 1}, // The buckets is for the histogram.
	})
	return &MetricsMiddleware{
		next: next,
		reqCounterAgg: reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg: reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}


func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
		"took": time.Since(start),
		"err": err,
		"func": "AggregateDistance",
	}).Info("Aggregate Distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}


func (m *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
		"took": time.Since(start),
		"err": err,
		"obuID": obuID,
		"totalDist": distance,
		"totalAmount": amount,
	}).Info("CalculateInvoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}
