package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	memoryCartItemsTotalCounter = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "cart",
			Name:      "memory_cart_items_total_counter",
			Help:      "Total number of items in the in-memory cart storage",
		},
	)

	httpRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "http_requests_total_counter",
			Help:      "Total number of HTTP requests received by the service, categorized by handler.",
		}, []string{"handler"},
	)

	externalRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "external_requests_total_counter",
			Help:      "Total number of requests to external resources, categorized by target service and handler.",
		}, []string{"service", "handler"},
	)

	httpResponseStatusTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "http_response_status_total_counter",
			Help:      "Total number of HTTP request execution statuses to external resources, including handler and status_code.",
		}, []string{"handler", "status_code"},
	)

	externalResponseStatusTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "cart",
			Name:      "external_response_status_total_counter",
			Help:      "Total number of external request execution statuses to external resources, including handler and status_code.",
		}, []string{"handler", "status_code"},
	)

	externalRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "cart",
			Name:      "external_requests_duration_histogram",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"service", "handler"})

	httpRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "cart",
			Name:      "http_requests_duration_histogram",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler"})
)

func UpdateMemoryCartItemsTotalCounter(cartItemsCount int) {
	memoryCartItemsTotalCounter.Set(float64(cartItemsCount))
}

func IncHttpRequestsTotalCounter(labelValues ...string) {
	httpRequestsTotalCounter.WithLabelValues(labelValues...).Inc()
}

func IncExternalRequestsTotalCounter(labelValues ...string) {
	externalRequestsTotalCounter.WithLabelValues(labelValues...).Inc()
}

func ObserveHttpRequestsDurationHistogram(createdAt time.Time, labelValues ...string) {
	httpRequestsDurationHistogram.WithLabelValues(labelValues...).Observe(time.Since(createdAt).Seconds())
}

func ObserveExternalRequestsDurationHistogram(createdAt time.Time, labelValues ...string) {
	externalRequestsDurationHistogram.WithLabelValues(labelValues...).Observe(time.Since(createdAt).Seconds())
}

func IncHttpResponseStatusTotalCounter(labelValues ...string) {
	httpResponseStatusTotalCounter.WithLabelValues(labelValues...).Inc()
}

func IncExternalResponseStatusTotalCounter(labelValues ...string) {
	externalResponseStatusTotalCounter.WithLabelValues(labelValues...).Inc()
}
