package relay

import "fmt"

type Device interface {
	Name() string
	On()
	Off()
	State() (bool, error)
	Number() uint8
}

type Signal uint8

const (
	Off Signal = iota
	On
)

type Service struct {
	fan    Device
	lights Device
}

func New(opts ...DeviceOpt) *Service {
	s := Service{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

type DeviceOpt func(s *Service)

func WithLightDevice(lights Device) DeviceOpt {
	return func(s *Service) {
		s.lights = lights
	}
}

func WithFanDevice(fan Device) DeviceOpt {
	return func(s *Service) {
		s.fan = fan
	}
}

func (s *Service) SetAllTo(signal Signal) {
	for _, device := range []Device{s.fan, s.lights} {
		switch signal {
		case Off:
			device.Off()
		case On:
			device.On()
		}
	}
}

func (s *Service) LoadedDevices() []Device {
	var devices []Device
	devices = append(devices, s.fan)
	devices = append(devices, s.lights)
	return devices
}

func (s *Service) DeviceByRelay(number uint8) (Device, error) {
	for _, device := range s.LoadedDevices() {
		if device.Number() == number {
			return device, nil
		}
	}
	return nil, fmt.Errorf("unable to find relay on that number")
}

func (s *Service) Fans() Device {
	return s.fan
}

func (s *Service) Lights() Device {
	return s.lights
}
