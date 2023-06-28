package controller

import (
	"github.com/hashicorp/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

var (
	lastLightStateOff time.Time
	lastLightStateOn  time.Time
)

func (c *Service) LightsControl() {

	lights := c.Relays.Lights()

	onTime := viper.GetTime("lights.onTime")
	now := time.Now()
	onTime = time.Date(now.Year(), now.Month(), now.Day(), onTime.Hour(), onTime.Minute(), onTime.Second(), onTime.Nanosecond(), time.Local)
	onDuration := viper.GetDuration("lights.Duration")
	offTime := onTime.Add(onDuration)
	prevOffTime := offTime.Add(time.Hour * -24)
	log.Debug().Msgf("On Time: %s, duration: %s, offTime: %s previous offTime: %s", onTime.String(), onDuration.String(), offTime.String(), prevOffTime.String())
	lightsOn, err := lights.State()
	if err != nil {
		metrics.IncrCounter([]string{"relay_light_state_error"}, 1)
		log.Error().Msgf("unable to get lights state: %s", err)
		return
	}

	if lastLightStateOff.IsZero() && !lightsOn {
		lastLightStateOff = time.Now()
	} else if lastLightStateOn.IsZero() && lightsOn {
		lastLightStateOn = time.Now()
	}

	if now.After(onTime) && now.Before(offTime) {
		if lightsOn {
			log.Info().Msg("Lights are already on, we are between on and off time not doing anything")
		} else {
			log.Info().Msg("Turning on lights, we are between the on and off time")
			lights.On()
			lightsOn = true
			lastLightStateOn = time.Now()
		}
	}

	if now.After(offTime) {
		if lightsOn {
			log.Info().Msg("Turning off lights, we are after the off time")
			lights.Off()
			lightsOn = false
			lastLightStateOff = time.Now()
		} else {
			log.Info().Msg("Lights are already off, we are after the off time")
		}
	}

	if now.Before(onTime) && now.After(prevOffTime) {
		if lightsOn {
			log.Info().Msg("Turning off lights, we are before the on time, and after the previous off time")
			lights.Off()
			lightsOn = false
			lastLightStateOn = time.Now()
		} else {
			log.Info().Msg("Lights are already off, we are before the on time, and after the previous off time")
		}
	}

	emitLightStateMetric(lightsOn)
}

func emitLightStateMetric(lightsOn bool) {

	var state float32
	if lightsOn {
		state = 1
	} else {
		state = 0
	}
	metrics.SetGauge([]string{"relay_light_state"}, state)
	if !lastLightStateOff.IsZero() {
		metrics.SetGauge([]string{"last_light_change_off_timestamp_seconds"}, float32(lastLightStateOff.Unix()))
	}
	if !lastLightStateOn.IsZero() {
		metrics.SetGauge([]string{"last_light_change_on_timestamp_seconds"}, float32(lastLightStateOn.Unix()))
	}
}
