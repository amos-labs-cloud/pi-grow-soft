package controller

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/hashicorp/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

var waterNeeded bool
var waterNeededTimeStamp time.Time

func (c *Service) WaterControl() {
	soilSensor := c.Sensors.MoistureSensor()

	capacitance, err := soilSensor.ReadMoisture(10)
	if err != nil {
		metrics.IncrCounter([]string{"moisture_read_error"}, 1)
		log.Warn().Msgf("unable to read moisture: %s", err)
	}
	log.Info().Msgf("capacitance: %.0f", capacitance)

	prevWaterNeeded := waterNeeded
	lowRange := viper.GetInt("moistureSensor.lowRange")
	waterNeeded = lowRange > int(capacitance)
	if waterNeeded {
		log.Info().Msgf("Water is needed capacitance is below threshold: %d current: %d", lowRange, int(capacitance))
		// If water needed was set to false, we are about set it to true, so set the new timestamp
		if !prevWaterNeeded {
			log.Debug().Msg("Setting water timestamp")
			waterNeededTimeStamp = time.Now()
		}
		log.Debug().Msg("setting water needed to true")
		metrics.SetGauge([]string{"last_water_needed_timestamp_seconds"}, float32(waterNeededTimeStamp.Unix()))
	} else {
		log.Debug().Msg("waterNeeded is false")
	}

	emitCapacitanceMetric(capacitance)
}

func emitCapacitanceMetric(capacitance float32) {
	metrics.SetGauge([]string{"moisture_capacitance"}, capacitance)
}

func (c *Service) WaterDisplayView() display.Page {
	var pages []string

	if waterNeeded {
		pages = append(pages, fmt.Sprintf("Water Me!!!"))
		return pages
	} else {
		pages = append(pages, fmt.Sprintf("Water is OK"))
		return pages
	}
}
