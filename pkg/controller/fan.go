package controller

import (
	"fmt"
	"github.com/hashicorp/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"math"
	"time"
)

var (
	lastFanStateOff time.Time
	lastFanStateOn  time.Time
)

func (c *Service) FanControl() {
	temp, _, err := c.readTempHumidity()
	if err != nil {
		log.Warn().Msgf("unable to read air temp sensor: %s", err)
		return
	}

	triggerFanTemp := viper.GetInt("fans.triggerTempCelsius")
	fans := c.Relays.Fans()
~.	roundedTemp := int(math.Round(float64(temp)))
	fansOn, err := fans.State()
	if err != nil {
		log.Error().Msgf("unable to get fan state: %s", err)
		return
	}
	if lastFanStateOff.IsZero() && !fansOn {
		lastLightStateOff = time.Now()
	} else if lastFanStateOn.IsZero() && fansOn {
		lastLightStateOn = time.Now()
	}

	if roundedTemp > triggerFanTemp {
		if fansOn {
			log.Info().Msgf("current temp: %d is greater than %d. fans already on, waiting for temp to lower", roundedTemp, triggerFanTemp)
		} else {
			log.Info().Msgf("current temp: %d is greater than %d, turning on the fan", roundedTemp, triggerFanTemp)
			fans.On()
			fansOn = true
			lastFanStateOn = time.Now()
		}
	}

	if roundedTemp <= triggerFanTemp {
		if fansOn {
			log.Info().Msgf("current temp: %d is under %d, turning fan off", roundedTemp, triggerFanTemp)
			fans.Off()
			fansOn = false
			lastFanStateOff = time.Now()
		} else {
			log.Info().Msgf("current temp: %d is under %d, not triggering fans", roundedTemp, triggerFanTemp)
		}
	}
	emitFanStateMetric(fansOn)
}

func emitFanStateMetric(fansOn bool) {

	var state float32
	if fansOn {
		state = 1
	} else {
		state = 0
	}
	metrics.SetGauge([]string{"relay_fan_state"}, state)
	if !lastFanStateOff.IsZero() {
		metrics.SetGauge([]string{"last_fan_change_off_timestamp_seconds"}, float32(lastFanStateOff.Unix()))
	}
	if !lastFanStateOn.IsZero() {
		metrics.SetGauge([]string{"last_fan_change_on_timestamp_seconds"}, float32(lastFanStateOn.Unix()))
	}
}

func (c *Service) readTempHumidity() (float32, float32, error) {
	airSensor, err := c.Sensors.AirTempSensor()
	if err != nil {
		// This is not an interesting metric because we are retrieving internal interfaces to then actually do something
		log.Error().Msgf("unable to get air sensor: %s", err)
		return 0, 0, fmt.Errorf("")
	}
	tries := 10
	temp, humidity, err := airSensor.ReadTempHumidity(tries)
	if err != nil {
		// This is interesting, can maybe do a method that logs a message, gives you an error to return with that
		// previous message, and submits a metric
		c.Metrics.IncrCounter([]string{"air_read_error"}, 1)

		log.Error().Msgf("unable to read temp after %d tries, try sudo", tries)
		return 0, 0, nil
	}
	log.Info().Msgf("temp: %.0f Celsius humidity: %.0f%%", temp, humidity)
	c.Metrics.SetGauge([]string{"air_temperature"}, temp)
	c.Metrics.SetGauge([]string{"air_humidity"}, humidity)
	return temp, humidity, nil
}
