package measurement

import (
	"fmt"
	"github.com/hashicorp/go-metrics"
	"github.com/hashicorp/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

var noWebServer bool

type Service struct {
	hostname    string
	serviceName string
	sinkName    string
	*metrics.Metrics
}

func WithNoWebServer() ConfigOpt {
	return func(s *Service) {
		noWebServer = true
	}
}

func WithServiceName(serviceName string) ConfigOpt {
	return func(s *Service) {
		s.serviceName = serviceName
	}
}

func WithSinkName(sinkName string) ConfigOpt {
	return func(s *Service) {
		s.sinkName = sinkName
	}
}

type ConfigOpt func(*Service)

func New(opts ...ConfigOpt) (*Service, error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Warn().Msgf("unable to get hostname for measurement")
	}

	s := Service{hostname: hostname}
	for _, opt := range opts {
		opt(&s)
	}

	sink, err := prometheus.NewPrometheusSinkFrom(prometheus.PrometheusOpts{
		Name: s.sinkName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create prometheus sink: %w", err)
	}

	m, err := metrics.New(&metrics.Config{
		ServiceName:          s.serviceName,
		EnableHostname:       true,
		EnableHostnameLabel:  true,
		EnableServiceLabel:   true,
		EnableRuntimeMetrics: true,
		EnableTypePrefix:     true,
		TimerGranularity:     time.Millisecond,
		ProfileInterval:      time.Second,
		FilterDefault:        true,
	}, sink)

	if err != nil {
		return nil, err
	}

	m.EnableHostnameLabel = true
	s.Metrics = m
	if !noWebServer {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			if err = http.ListenAndServe(":8080", nil); err != nil {
				log.Error().Msgf("unable to start prometheus endpoint: %s", err)
			}
		}()
	}

	return &s, nil
}
