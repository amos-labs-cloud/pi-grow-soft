package controller

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/relay"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
)

func TestFanMetricsEmittance(t *testing.T) {

	metricsService, err := measurement.New(measurement.WithNoWebServer(), measurement.WithServiceName("fantest"), measurement.WithSinkName("fantest"))
	if err != nil {
		t.Fatal(err)
	}

	device := new(MockDevice)
	device.Mock.On("State").Return(func() (bool, error) { return true, nil })

	relayService := relay.New(relay.WithFanDevice(device))

	dht11 := new(MockDHT11)
	dht11.Mock.On("ReadTempHumidity", mock.Anything).Return(
		func() (float32, float32, error) {
			min := float32(23)
			max := float32(35)
			return min + rand.Float32()*(max-min), min + rand.Float32()*(max-min), nil
		})

	sensorService := sensors.New(sensors.WithAirSensor(dht11))
	controller := New(
		WithMetricsService(metricsService),
		WithSensorService(sensorService),
		WithRelayService(relayService),
	)

	for i := 0; i < 2; i++ {
		controller.FanControl()
		time.Sleep(1 * time.Second)
	}

}
