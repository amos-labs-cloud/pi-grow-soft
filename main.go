package main

import (
	"context"
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/cmd"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/controller"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/pin"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/relay"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/moisture"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/temp_humidity"
	"github.com/d2r2/go-logger"
	lcd1602 "github.com/pimvanhespen/go-pi-lcd1602"
	"github.com/pimvanhespen/go-pi-lcd1602/synchronized"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"os"
	"strings"
)

func main() {
	if err := loadConfig(); err != nil {
		panic(err)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	err := logger.ChangePackageLogLevel("dht", logger.FatalLevel)
	if err != nil {
		log.Panic().Msgf("could not turn off dht logger: %s", err)
	}
	controllerService := setup()

	cmd.Initialize(controllerService)
	cmd.Execute()
}

func loadConfig() error {
	viper.SetDefault("relay.one", uint8(5))
	viper.SetDefault("relay.two", uint8(6))
	viper.SetDefault("relay.three", uint8(13))
	viper.SetDefault("relay.four", uint8(19))

	viper.SetConfigName("controller-config") // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}

func WSL() bool {
	magicString := "microsoft"
	b, err := os.ReadFile("/proc/version")
	if err != nil {
		panic(err)
	}

	return strings.Contains(string(b), magicString)
}

func setup() *controller.Service {
	var controllerService *controller.Service
	if WSL() {
		lights := new(controller.MockDevice)
		lights.Mock.On("Off").Return()
		lights.Mock.On("On").Return()

		fans := new(controller.MockDevice)
		fans.Mock.On("Off").Return()
		fans.Mock.On("On").Return()

		dht11 := new(controller.MockDHT11)
		dht11.Mock.On("ReadTempHumidity", mock.Anything).Return(
			func() (float32, float32, error) {
				min := float32(20)
				max := float32(28)
				return min + rand.Float32()*(max-min), min + rand.Float32()*(max-min), nil
			})
		fans.Mock.On("State").Return(func() (bool, error) { return true, nil })
		lights.Mock.On("State").Return(func() (bool, error) { return false, nil })
		soilSensor := new(controller.MockDevice)
		soilSensor.Mock.On("ReadMoisture", mock.Anything).Return(func() float32 {
			min := float32(1)
			max := float32(1000)
			return min + rand.Float32()*(max-min)
		}, nil)
		relayService := relay.New(relay.WithFanDevice(fans), relay.WithLightDevice(lights))
		sensorService := sensors.New(sensors.WithAirSensor(dht11), sensors.WithMoistureSensor(soilSensor))
		metricsService, err := measurement.New()
		if err != nil {
			log.Panic().Msgf("could not create metricService: %s", err)
		}

		controllerService = controller.New(
			controller.WithSensorService(sensorService),
			controller.WithRelayService(relayService),
			controller.WithMetricsService(metricsService),
		)
	} else {
		log.Info().Msg("Creating relays")

		lights := relay.NewRelay("lights", uint8(viper.GetInt("relay.one")), uint8(1), relay.NormallyOpen)
		fan := relay.NewRelay("fan", uint8(viper.GetInt("relay.two")), uint8(2), relay.NormallyOpen)
		relayService := relay.New(relay.WithFanDevice(fan), relay.WithLightDevice(lights))
		airTHSensor := temp_humidity.NewTHSensor("air", uint8(viper.GetInt("thSensor.pin")), temp_humidity.DHT11)
		soilSensor := moisture.NewMoistureSensor("soil", 0x36, "/dev/i2c-1")
		sensorService := sensors.New(sensors.WithAirSensor(airTHSensor), sensors.WithMoistureSensor(soilSensor))
		metricsService, err := measurement.New()
		if err != nil {
			log.Panic().Msgf("could not create metric service: %s", err)
		}
		pinService := pin.New()
		lcdi := lcd1602.New(
			viper.GetInt("display.rs"),
			viper.GetInt("display.e"),
			viper.GetIntSlice("display.pins"),
			viper.GetInt("display.width"),
		)
		syncedLCD := synchronized.NewSynchronizedLCD(lcdi)
		syncedLCD.Initialize()
		defer syncedLCD.Close()
		displayService := display.New(syncedLCD)

		ctx := context.Background()
		displayService.AddPage("PiGrow", "by Amos Labs")

		go func() {
			displayService.RunPages(ctx)
		}()

		controllerService = controller.New(
			controller.WithDisplayService(displayService),
			controller.WithPinService(pinService),
			controller.WithSensorService(sensorService),
			controller.WithRelayService(relayService),
			controller.WithMetricsService(metricsService),
		)
		displayService.AddPageFunc(airTHSensor.DisplayView)
		displayService.AddPageFunc(soilSensor.DisplayView)
		displayService.AddPageFunc(controllerService.WaterDisplayView)
	}
	return controllerService
}
