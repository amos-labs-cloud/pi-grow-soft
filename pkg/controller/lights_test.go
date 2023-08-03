package controller

import (
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/relay"
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestService_LightsControl(t *testing.T) {
	viper.SetConfigName("controller-config") // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../..")             // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		t.Fatalf("fatal error config file: %s", err)
	}

	t.Logf("%s", time.Now().String())
	onTime := viper.GetTime("lights.onTime")
	onTime = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), onTime.Hour(), onTime.Minute(), onTime.Second(), onTime.Nanosecond(), time.Local)
	onDuration := viper.GetDuration("lights.Duration")
	offTime := onTime.Add(onDuration)
	t.Logf("On Time: %s, duration: %s, offTime: %s", onTime.String(), onDuration.String(), offTime.String())
	if time.Now().After(onTime) && time.Now().Before(offTime) {
		t.Log("Turning on lights, we are between the on and off time")
	}

	if time.Now().After(offTime) {
		t.Log("Turning off lights, we are after the off time")
	}

	if time.Now().Before(onTime) {
		t.Log("Turning off lights, we are before the on time")
	}
}

func TestLightMetricsEmittance(t *testing.T) {
	metricsService, err := measurement.New(measurement.WithNoWebServer(), measurement.WithServiceName("lighttest"), measurement.WithSinkName("lighttest"))
	if err != nil {
		t.Fatal(err)
	}

	device := new(MockDevice)
	//device.Mock.On("State").Return(func() (bool, error) { return rand.Intn(2) == 1, nil })
	//device.Mock.On("State").Return(func() (bool, error) { return true, nil })
	device.Mock.On("Off").Return()
	device.Mock.On("Off").Return()
	device.Mock.On("On").Return()

	relayService := relay.New(relay.WithDevice(device))

	controller := New(
		WithMetricsService(metricsService),
		WithRelayService(relayService),
	)

	for i := 0; i < 2; i++ {
		viper.Set("lights.OnTime", "10:00AM")
		viper.Set("lights.Duration", "3h")
		device.Mock.On("State").Return(func() (bool, error) { return false, nil })
		controller.LightsControl()
		time.Sleep(1 * time.Second)
		viper.Set("lights.OnTime", "12:00PM")
		viper.Set("lights.Duration", "3h")
		controller.LightsControl()
	}

}
