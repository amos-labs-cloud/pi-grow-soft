package relay

import (
	"fmt"
)

type Device interface {
	Name() string
	On()
	Off()
	State() (bool, error)
	Number() int
	TypeInfo() DeviceTypeInfo
}

type DeviceTypeInfo struct {
	Category DeviceCategory
}

type DeviceCategory int

const (
	Fans DeviceCategory = iota
	Lights
)

func (d DeviceCategory) String() string {
	switch d {
	case Fans:
		return "fans"
	case Lights:
		return "lights"
	}
	return "unknown"
}

type Signal uint8

const (
	Off Signal = iota
	On
)

type Service struct {
	devices map[string]Device
}

func New(opts ...DeviceOpt) *Service {
	devices := make(map[string]Device)
	s := Service{devices: devices}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

type DeviceOpt func(s *Service)

func WithDevice(device Device) DeviceOpt {
	return func(s *Service) {
		s.devices[device.TypeInfo().Category.String()] = device
	}
}

func (s *Service) SetAllTo(signal Signal) {
	for _, device := range s.devices {
		switch signal {
		case Off:
			device.Off()
		case On:
			device.On()
		}
	}
}

func (s *Service) Devices() map[string]Device {
	return s.devices
}

func (s *Service) DeviceByRelay(number int) (Device, error) {
	for _, device := range s.Devices() {
		if device.Number() == number {
			return device, nil
		}
	}
	return nil, fmt.Errorf("unable to find relay on that number")
}

func (s *Service) Fans() Device {
	return s.devices[Fans.String()]
}

func (s *Service) Lights() Device {
	return s.devices[Lights.String()]
}
