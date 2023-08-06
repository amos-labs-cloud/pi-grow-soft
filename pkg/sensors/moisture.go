package sensors

type MoistureSensorType uint

const (
	AdaSoilSensor MoistureSensorType = iota
)

type Config struct {
	Enabled         bool                   `mapstructure:"enabled"`
	DefaultLowRange int                    `mapstructure:"defaultLowRange"`
	SensorMapping   map[string]interface{} `mapstructure:"sensorMapping"`
}

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
