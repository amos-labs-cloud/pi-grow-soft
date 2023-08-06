package main

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/cmd"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/controller"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/measurement"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/pin"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/relay"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/moisture"
	dht11 "github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors/temp_humidity"
	"github.com/d2r2/go-logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
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
	viper.SetDefault("relay.one", uint8(6))
	viper.SetDefault("relay.two", uint8(5))

	viper.SetConfigName("controller-config") // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}

// TODO: There needs to be a better way to run tests off-Pi
func setup() *controller.Service {
	var controllerService *controller.Service
	relayService := relay.New(relaysToLoad()...)
	sensorService := sensors.New(sensorsToLoad()...)

	metricsService, err := measurement.New()
	if err != nil {
		log.Panic().Msgf("could not create metric service: %s", err)
	}
	pinService := pin.New()
	controllerService = controller.New(
		controller.WithPinService(pinService),
		controller.WithSensorService(sensorService),
		controller.WithRelayService(relayService),
		controller.WithMetricsService(metricsService),
	)

	return controllerService
}

func sensorsToLoad() []sensors.SensorOpt {
	var toLoad []sensors.SensorOpt
	if viper.GetBool("thSensor.enabled") {
		thSensorPin := viper.GetInt("thSensor.pin")
		log.Info().Msgf("loading air temp humidity sensor on pin: %d", thSensorPin)
		airTHSensor := dht11.NewDHT11(thSensorPin)
		toLoad = append(toLoad, sensors.WithSensor(airTHSensor, 1))
	}

	if viper.GetBool("moistureSensors.enabled") {
		log.Info().Msgf("loading soil moisture sensor(s)")
		for i := 1; i < 4; i++ { // TODO: we can only have three sensors at the moment, hard code to move forward
			sensorConfig := viper.Sub(fmt.Sprintf("moistureSensors.%d", i))
			if sensorConfig == nil {
				log.Info().Msgf("not loading moisture sensor: %d, unable to get subkey config", i)
				continue
			}
			var sc moisture.StemmaConfig
			err := sensorConfig.Unmarshal(&sc)
			if err != nil {
				log.Info().Msgf("not loading moisture sensor: %d, unable to get read config in config struct: %s", i, err)
			}
			log.Info().Msgf("loading soil sensor at address %d", sc.Address)
			soilSensor := moisture.NewAdafruitStemma("soil", "/dev/i2c-1", sc)
			toLoad = append(toLoad, sensors.WithSensor(soilSensor, i))
		}
	}
	return toLoad
}

func relaysToLoad() []relay.DeviceOpt {
	var toLoad []relay.DeviceOpt
	var lights *relay.Relay
	if viper.GetBool("devices.lights.enabled") {
		lightsPin := viper.GetInt(fmt.Sprintf("relay.%d", viper.GetInt("devices.lights.relay")))
		log.Info().Msgf("creating relay for lights to be controlled on pin %d", lightsPin)
		lights = relay.NewRelay("lights", lightsPin, viper.GetInt("devices.lights.relay"), relay.NormallyOpen, relay.Lights)
		toLoad = append(toLoad, relay.WithDevice(lights))
	}

	var fans *relay.Relay
	if viper.GetBool("devices.fans.enabled") {
		fansPin := viper.GetInt(fmt.Sprintf("relay.%d", viper.GetInt("devices.fans.relay")))
		log.Info().Msgf("creating relay for fans to be controlled on pin %d", fansPin)
		fans = relay.NewRelay("fans", fansPin, viper.GetInt("devices.fans.relay"), relay.NormallyOpen, relay.Fans)
		toLoad = append(toLoad, relay.WithDevice(fans))
	}
	return toLoad
}
