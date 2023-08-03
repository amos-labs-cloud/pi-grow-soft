package sensors

type Service struct {
	sensors map[string]Sensor
}

type Sensor interface {
	Name() string
	SensorTypeInfo() SensorTypeInfo
}

type SensorTypeInfo struct {
	Category SensorCategory
	Info     map[string]string
}

type SensorCategory int

const (
	AirTempHumidity SensorCategory = iota
	Moisture
)

func (sc SensorCategory) String() string {
	switch sc {
	case AirTempHumidity:
		return "air_temp_humidity"
	case Moisture:
		return "moisture"
	}
	return "unknown"
}

type SensorOpt func(s *Service)

func WithSensor(sensor Sensor) SensorOpt {
	return func(s *Service) {
		category := sensor.SensorTypeInfo().Category.String()
		s.sensors[category] = sensor
	}
}

func New(opts ...SensorOpt) *Service {
	s := Service{}
	s.sensors = make(map[string]Sensor)
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func (s *Service) AirTempHumiditySensor() TempHumiditySensor {
	return s.sensors[AirTempHumidity.String()].(TempHumiditySensor)
}

func (s *Service) MoistureSensor() MoistureSensor {
	return s.sensors[Moisture.String()].(MoistureSensor)
}
