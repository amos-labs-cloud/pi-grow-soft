package temp_humidity

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/sensors"
	"github.com/d2r2/go-dht"
	"github.com/stianeikeland/go-rpio/v4"
)

type DHT11 struct {
	pinNumber int
	tries     int
}

func NewDHT11(pinNumber int) *DHT11 {
	return &DHT11{pinNumber: pinNumber, tries: 10}
}

func (ths *DHT11) ReadTempHumidity() (float32, float32, error) {
	err := rpio.Open()
	if err != nil {
		return 0, 0, fmt.Errorf("unable to open rpio: %s", err)
	}
	defer rpio.Close()

	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, int(ths.pinNumber), false, ths.tries)
	if err != nil {
		return 0, 0, fmt.Errorf("ran into error ready sensor: %+v err: %w", ths, err)
	}
	return temperature, humidity, nil
}

func (ths *DHT11) Name() string {
	return "dht11_temp_humidity_sensor"
}

func (ths *DHT11) SensorTypeInfo() sensors.SensorTypeInfo {
	return sensors.SensorTypeInfo{
		Category: sensors.AirTempHumidity,
		Info: map[string]string{
			"description": "dht11 digital temp and humidity sensor",
		},
	}
}

//func (ths *DHT11) DisplayView() display.Page {
//	var pages display.Page
//	log.Debug().Msgf("displayView thSensor: %+v", ths)
//	if ths.lastTemp != nil && ths.lastHumidity != nil {
//		roundedTemp := int(math.Round(float64(*ths.lastTemp)))
//		tempMessage := fmt.Sprintf("Air Temp: %dßCC", roundedTemp)
//		humidityMessage := fmt.Sprintf("Humidity: %.0f%%", *ths.lastHumidity)
//		pages = append(pages, tempMessage)
//		pages = append(pages, humidityMessage)
//		return pages
//	}
//	notReady := "TempHumiditySensor != ready"
//	pages = append(pages, notReady)
//	return pages
//}
