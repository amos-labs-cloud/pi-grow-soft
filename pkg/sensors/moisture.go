package sensors

type MoistureSensorType uint

const (
	AdaSoilSensor MoistureSensorType = iota
)

func (m MoistureSensorType) String() string {
	switch m {
	case AdaSoilSensor:
		return "adafruit_moisture_sensor"
	}
	return "unknown"
}

type MoistureSensor interface {
	ReadMoisture(tries int) (float32, error)
}
