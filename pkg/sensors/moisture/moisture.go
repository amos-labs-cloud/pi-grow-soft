package moisture

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/display"
	"github.com/asssaf/stemma-soil-go/soil"
	"github.com/rs/zerolog/log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type MoistureSensorType uint

const (
	AdaSoilSensor MoistureSensorType = iota
)

type MoistureDevice struct {
	lastCapacitance *float32
	label           string
	dev             *soil.Dev
	closeFunc       func() error
	devAddr         uint16
	devPath         string
}

type Sensor interface {
	ReadMoisture(tries int) (float32, error)
}

func NewMoistureSensor(label string, devAddr uint16, devPath string) *MoistureDevice {
	return &MoistureDevice{label: label, devAddr: devAddr, devPath: devPath}
}

func (m *MoistureDevice) ReadMoisture(tries int) (float32, error) {
	m.start()
	values := soil.SensorValues{}
	if err := m.dev.Sense(&values); err != nil {
		return 0, fmt.Errorf("unable to sense moisture: %w", err)
	}

	floatCap := float32(values.Capacitance)
	m.lastCapacitance = &floatCap
	m.stop()
	return *m.lastCapacitance, nil
}

func (m *MoistureDevice) start() {
	if _, err := host.Init(); err != nil {
		panic(err)
	}

	i2cPort, err := i2creg.Open(m.devPath)
	if err != nil {
		panic(err)
	}

	opts := soil.DefaultOpts
	if m.devAddr != 0 {
		if m.devAddr < 0x36 || m.devAddr > 0x39 {
			panic(fmt.Sprintf("given address not supported by device: %x", m.devAddr))
		}
		opts.Addr = m.devAddr
	}

	dev, err := soil.NewI2C(i2cPort, &opts)
	if err != nil {
		panic(err)
	}

	m.closeFunc = dev.Halt
	m.dev = dev
}

func (m *MoistureDevice) stop() {
	err := m.closeFunc()
	if err != nil {
		log.Error().Msgf("unable to close rpio: %s", err)
	}
}

func (m *MoistureDevice) DisplayView() display.Page {
	var pages []string
	log.Debug().Msgf("displayView: moisture sensor: %+v", m)
	if m.lastCapacitance != nil {
		capacitanceMessage := fmt.Sprintf("Moisture: %.0f", *m.lastCapacitance)
		pages = append(pages, capacitanceMessage)
		return pages
	}
	notReady := "Moisture != ready"
	pages = append(pages, notReady)
	return pages
}
