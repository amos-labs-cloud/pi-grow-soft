package sensors

type Service struct {
	sensors map[SensorCategory]map[int]Sensor
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

func WithSensor(sensor Sensor, number int) SensorOpt {
	return func(s *Service) {
		category := sensor.SensorTypeInfo().Category
		if s.sensors[category] == nil {
			s.sensors[category] = make(map[int]Sensor)
		}
		s.sensors[category][number] = sensor
	}
}

func New(opts ...SensorOpt) *Service {
	s := Service{}
	s.sensors = make(map[SensorCategory]map[int]Sensor)
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func (s *Service) AirTempHumiditySensors() map[int]TempHumiditySensor {
	retMap := map[int]TempHumiditySensor{}
	sensors := s.getSensorByCategory(AirTempHumidity)
	for i, v := range sensors {
		retMap[i] = v.(TempHumiditySensor)
	}
	return retMap
}

func (s *Service) MoistureSensors() map[int]MoistureSensor {
	retMap := map[int]MoistureSensor{}
	sensors := s.getSensorByCategory(Moisture)
	for i, v := range sensors {
		retMap[i] = v.(MoistureSensor)
	}
	return retMap
}

func (s *Service) getSensorByCategory(cat SensorCategory) map[int]interface{} {
	sensors := s.sensors[cat]
	var retMap map[int]interface{}
	retMap = make(map[int]interface{})
	for i, sensor := range sensors {
		switch cat {
		case AirTempHumidity:
			retMap[i] = sensor.(TempHumiditySensor)
		case Moisture:
			retMap[i] = sensor.(MoistureSensor)
		}
	}
	return retMap
}
