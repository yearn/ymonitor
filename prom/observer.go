package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tcnksm/go-httpstat"
)

var labels = []string{"host", "network", "env", "code", "type"}

const NS = "ymonitor"
const SUB = "node"

var requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "request_counter",
	Help:      "A simple request counter",
}, labels)
var durationBuckets = prometheus.LinearBuckets(0.025, 0.025, 400)
var connectDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "connect_time_seconds",
	Help:      "Histogram of connect_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var dnsLookupDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "dns_lookup_time_seconds",
	Help:      "Histogram of dns_lookup_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var nameLookupDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "name_lookup_time_seconds",
	Help:      "Histogram of name_lookup_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var preTransferDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "pre_transfer_time_seconds",
	Help:      "Histogram of pre_transfer_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var startTransferDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "start_transfer_time_seconds",
	Help:      "Histogram of start_transfer_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var tcpConnectionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "tcp_connection_time_seconds",
	Help:      "Histogram of tcp_connection_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var tlsHandshakeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "tls_handshake_time_seconds",
	Help:      "Histogram of tls_handshake_time in seconds.",
	Buckets:   durationBuckets,
}, labels)
var serverProcessingDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: NS,
	Subsystem: SUB,
	Name:      "server_processing_time_seconds",
	Help:      "Histogram of server_processing_time in seconds.",
	Buckets:   durationBuckets,
}, labels)

func Observe(stats httpstat.Result, labels prometheus.Labels) {
	requestCounter.With(labels).Inc()
	connectDuration.With(labels).Observe(float64(stats.Connect.Seconds()))
	dnsLookupDuration.With(labels).Observe(float64(stats.DNSLookup.Seconds()))
	nameLookupDuration.With(labels).Observe(float64(stats.NameLookup.Seconds()))
	preTransferDuration.With(labels).Observe(float64(stats.Pretransfer.Seconds()))
	startTransferDuration.With(labels).Observe(float64(stats.StartTransfer.Seconds()))
	tcpConnectionDuration.With(labels).Observe(float64(stats.TCPConnection.Seconds()))
	tlsHandshakeDuration.With(labels).Observe(float64(stats.TLSHandshake.Seconds()))
	serverProcessingDuration.With(labels).Observe(float64(stats.ServerProcessing.Seconds()))
}
