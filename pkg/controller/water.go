package controller

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/hashicorp/go-metrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var waterNeeded bool
var waterNeededTimeStamp time.Time
var moistureTracker map[int]moistureInfo

type moistureInfo struct {
	waterNeeded          bool
	waterNeededTimeStamp time.Time
	sensor               sensors.MoistureSensor
}

func (mi moistureInfo) isEmpty() bool {
	if mi.waterNeeded == false && mi.waterNeededTimeStamp.IsZero() && mi.sensor == nil {
		return true
	}
	return false
}

func (c *Service) WaterControl() {
	mSensors := c.Sensors.MoistureSensors()
	log.Debug().Msgf("controller: loaded the following moisture sensors: %+v", mSensors)
	// Init the tracker if it hasn't been yet
	if moistureTracker == nil {
		moistureTracker = make(map[int]moistureInfo)
	}

	// for each sensor, create a tracker entry
	for i, soilSensor := range mSensors {
		if moistureTracker[i].isEmpty() {
			moistureTracker[i] = moistureInfo{sensor: soilSensor}
		}
	}

	log.Debug().Msgf("moistureTracker: %+v", moistureTracker)

	for i, mInfo := range moistureTracker {
		capacitance, err := mInfo.sensor.ReadMoisture(10)
		if err != nil {
			metrics.IncrCounter([]string{"moisture_read_error_" + strconv.Itoa(i)}, 1)
			log.Warn().Msgf("unable to read moisture on sensor: %d, err: %s", i, err)
		}
		log.Info().Msgf("sensor %d: capacitance: %.0f", i, capacitance)

		prevWaterNeeded := mInfo.waterNeeded
		lowRange := viper.GetInt("controller.moisture.lowRange")
		waterNeeded = lowRange > int(capacitance)
		if waterNeeded {
			log.Info().Msgf("sensor: %d, water is needed capacitance is below threshold: %d current: %d", i, lowRange, int(capacitance))
			// If water needed was set to false, we are about set it to true, so set the new timestamp
			if !prevWaterNeeded {
				log.Debug().Msgf("sensor: %d Setting water timestamp", i)
				waterNeededTimeStamp = time.Now()
			}
			log.Debug().Msgf("sensor: %d setting water needed to true", i)
			metrics.SetGauge([]string{"last_water_needed_timestamp_seconds_" + strconv.Itoa(i)}, float32(waterNeededTimeStamp.Unix()))
		} else {
			log.Debug().Msgf("sensor: %d waterNeeded is false", i)
		}

		emitCapacitanceMetric(i, capacitance)
	}
}

func emitCapacitanceMetric(sensorNumber int, capacitance float32) {
	metrics.SetGauge([]string{"moisture_capacitance_" + strconv.Itoa(sensorNumber)}, capacitance)
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
