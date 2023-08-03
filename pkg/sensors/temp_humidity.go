package sensors

type THSensorModel uint

const (
	DHT11 THSensorModel = iota
)

type TempHumiditySensor interface {
	ReadTempHumidity() (float32, float32, error)
}
