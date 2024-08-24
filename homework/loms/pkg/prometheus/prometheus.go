package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	grpcRequestsTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "loms",
			Name:      "grpc_requests_total_counter",
			Help:      "Total number of gRPC requests received by the service, categorized by handler.",
		}, []string{"handler"},
	)

	dbRequestsTotalCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "loms",
		Name:      "db_requests_total_counter",
		Help:      "Total number of database requests, categorized by query type.",
	}, []string{"query_type"})

	grpcResponseStatusTotalCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "loms",
			Name:      "grpc_response_status_total_counter",
			Help:      "Total number of gRPC request execution statuses to external resources, including handler and status_code.",
		}, []string{"handler", "status_code"},
	)

	grpcRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "grpc_requests_duration_histogram",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler"})

	dbRequestsDurationHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "db_requests_duration_histogram",
			Help:      "Duration of database requests in seconds, categorized by query type and error status.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"query_type", "status"},
	)
)

func IncGRPCRequestsTotalCounter(labelValues ...string) {
	grpcRequestsTotalCounter.WithLabelValues(labelValues...).Inc()
}

func IncDBRequestsTotalCounter(labelValues ...string) {
	dbRequestsTotalCounter.WithLabelValues(labelValues...).Inc()
}

func ObserveGRPCRequestsDurationHistogram(createdAt time.Time, labelValues ...string) {
	grpcRequestsDurationHistogram.WithLabelValues(labelValues...).Observe(time.Since(createdAt).Seconds())
}

func IncGRPCResponseStatusTotalCounter(labelValues ...string) {
	grpcResponseStatusTotalCounter.WithLabelValues(labelValues...).Inc()
}

func ObserveDBRequestsDurationHistogram(startTime time.Time, labelValues ...string) {
	dbRequestsDurationHistogram.WithLabelValues(labelValues...).Observe(time.Since(startTime).Seconds())
}
