package controller

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
)

func TestWaterMetrics(t *testing.T) {
	rand.NewSource(time.Now().Unix())
	metricsService, err := measurement.New(measurement.WithNoWebServer(), measurement.WithServiceName("watertest"), measurement.WithSinkName("watertest"))
	if err != nil {
		t.Fatal(err)
	}

	device := new(MockDevice)
	device.Mock.On("SensorTypeInfo").Return(sensors.SensorTypeInfo{
		Category: sensors.Moisture,
	})

	sensorService := sensors.New(sensors.WithSensor(device))
	controller := New(
		WithMetricsService(metricsService),
		WithSensorService(sensorService),
	)

	for i := 0; i < 2; i++ {
		viper.Set("moistureSensor.lowRange", 300)
		device.Mock.On("ReadMoisture", mock.Anything).Return(func() float32 {
			min := float32(1)
			max := float32(1000)
			return min + rand.Float32()*(max-min)
		}, nil)
		controller.WaterControl()
		time.Sleep(1 * time.Second)
	}

}
