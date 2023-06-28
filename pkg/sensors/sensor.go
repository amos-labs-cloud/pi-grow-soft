package sensors

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/moisture"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/temp_humidity"
)

type Service struct {
	airTempHumiditySensor temp_humidity.Sensor
	moistureSensor        moisture.Sensor
}

type SensorOpt func(s *Service)

func WithAirSensor(sensor temp_humidity.Sensor) SensorOpt {
	return func(s *Service) {
		s.airTempHumiditySensor = sensor
	}
}

func WithMoistureSensor(sensor moisture.Sensor) SensorOpt {
	return func(s *Service) {
		s.moistureSensor = sensor
	}
}

func New(opts ...SensorOpt) *Service {
	s := Service{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func (s *Service) AirTempSensor() (temp_humidity.Sensor, error) {
	return s.airTempHumiditySensor, nil
}

func (s *Service) MoistureSensor() (moisture.Sensor, error) {
	return s.moistureSensor, nil
}
