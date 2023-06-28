package temp_humidity

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/d2r2/go-dht"
	"github.com/rs/zerolog/log"
	"math"
)

type THSensorType uint

const (
	DHT11 THSensorType = iota
)

type THDevice struct {
	sensorType   THSensorType
	closeFunc    func() error
	pinNumber    uint8
	label        string
	lastTemp     *float32
	lastHumidity *float32
}

type Sensor interface {
	ReadTempHumidity(tries int) (float32, float32, error)
}

func NewTHSensor(label string, pinNumber uint8, sensor THSensorType) *THDevice {
	return &THDevice{label: label, pinNumber: pinNumber, sensorType: sensor}
}

func (ths *THDevice) start() {}

func (ths *THDevice) stop() {}

func (ths *THDevice) ReadTempHumidity(tries int) (float32, float32, error) {
	ths.start()
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, int(ths.pinNumber), false, tries)
	if err != nil {
		return 0, 0, fmt.Errorf("ran into error ready sensor: %+v err: %w", ths, err)
	}
	ths.lastTemp = &temperature
	ths.lastHumidity = &humidity
	ths.stop()
	return temperature, humidity, nil
}

func (ths *THDevice) DisplayView() display.Page {
	var pages display.Page
	log.Debug().Msgf("displayView thSensor: %+v", ths)
	if ths.lastTemp != nil && ths.lastHumidity != nil {
		roundedTemp := int(math.Round(float64(*ths.lastTemp)))
		tempMessage := fmt.Sprintf("Air Temp: %dßCC", roundedTemp)
		humidityMessage := fmt.Sprintf("Humidity: %.0f%%", *ths.lastHumidity)
		pages = append(pages, tempMessage)
		pages = append(pages, humidityMessage)
		return pages
	}
	notReady := "THDevice != ready"
	pages = append(pages, notReady)
	return pages
}
