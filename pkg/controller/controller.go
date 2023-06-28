package controller

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/pin"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/relay"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
)

type Service struct {
	Display *display.Service
	Metrics *measurement.Service
	Pins    *pin.Service
	Relays  *relay.Service
	Sensors *sensors.Service
}

type Opts func(service *Service)

func WithDisplayService(displayService *display.Service) Opts {
	return func(s *Service) {
		s.Display = displayService
	}
}

func WithMetricsService(metricService *measurement.Service) Opts {
	return func(s *Service) {
		s.Metrics = metricService
	}
}

func WithPinService(pinService *pin.Service) Opts {
	return func(s *Service) {
		s.Pins = pinService
	}
}

func WithRelayService(relayService *relay.Service) Opts {
	return func(s *Service) {
		s.Relays = relayService
	}
}

func WithSensorService(sensorService *sensors.Service) Opts {
	return func(s *Service) {
		s.Sensors = sensorService
	}
}

func New(opts ...Opts) *Service {
	s := Service{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}
