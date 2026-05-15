package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	errorTypeServer     = "server"
	errorTypeNotFound   = "not_found"
	errorTypeValidation = "validation"
	unmatchedRouteLabel = "unmatched"
)

type Metrics struct {
	HTTP     HTTPMetrics
	Business BusinessMetrics
	Database DatabaseMetrics
}

type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
	ErrorsTotal      *prometheus.CounterVec
}

type BusinessMetrics struct {
	PostsCreatedTotal    prometheus.Counter
	PostsPublishedTotal  prometheus.Counter
	PostsDeletedTotal    prometheus.Counter
	CommentsCreatedTotal prometheus.Counter
	PostsTotal           *prometheus.GaugeVec
	CommentsTotal        prometheus.Gauge
}

type DatabaseMetrics struct {
	QueryDuration *prometheus.HistogramVec
	ErrorsTotal   *prometheus.CounterVec
}

func NewMetrics(registerer prometheus.Registerer) *Metrics {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}

	httpMetrics := HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of processed HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		RequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of in-flight HTTP requests",
			},
		),
		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total HTTP errors by category",
			},
			[]string{"type"},
		),
	}

	businessMetrics := BusinessMetrics{
		PostsCreatedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "posts_created_total",
				Help: "Total number of created posts",
			},
		),
		PostsPublishedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "posts_published_total",
				Help: "Total post transitions to published=true",
			},
		),
		PostsDeletedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "posts_deleted_total",
				Help: "Total number of deleted posts",
			},
		),
		CommentsCreatedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "comments_created_total",
				Help: "Total number of created comments",
			},
		),
		PostsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "posts_total",
				Help: "Current post count in database",
			},
			[]string{"published"},
		),
		CommentsTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "comments_total",
				Help: "Current comment count in database",
			},
		),
	}

	databaseMetrics := DatabaseMetrics{
		QueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_errors_total",
				Help: "Total number of database errors by operation",
			},
			[]string{"operation"},
		),
	}

	m := &Metrics{
		HTTP:     httpMetrics,
		Business: businessMetrics,
		Database: databaseMetrics,
	}

	mustRegister(
		registerer,
		m.HTTP.RequestsTotal,
		m.HTTP.RequestDuration,
		m.HTTP.RequestsInFlight,
		m.HTTP.ErrorsTotal,
		m.Business.PostsCreatedTotal,
		m.Business.PostsPublishedTotal,
		m.Business.PostsDeletedTotal,
		m.Business.CommentsCreatedTotal,
		m.Business.PostsTotal,
		m.Business.CommentsTotal,
		m.Database.QueryDuration,
		m.Database.ErrorsTotal,
	)

	return m
}

func (m *Metrics) RouteLabel(route string) string {
	if route == "" {
		return unmatchedRouteLabel
	}
	return route
}

func (m *Metrics) ErrorTypeForStatus(statusCode int) string {
	switch {
	case statusCode >= 500:
		return errorTypeServer
	case statusCode == 404:
		return errorTypeNotFound
	case statusCode >= 400:
		return errorTypeValidation
	default:
		return ""
	}
}

func mustRegister(registerer prometheus.Registerer, collectors ...prometheus.Collector) {
	for _, collector := range collectors {
		registerer.MustRegister(collector)
	}
}
