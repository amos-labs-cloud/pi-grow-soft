package relay

import (
	"fmt"
	"github.com/amos-labs-cloud/pi-grow-soft/pkg/pin"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

type Normality uint8

const (
	NormallyOpen Normality = iota
	NormallyClosed
	Undefined Normality = 8
)

func (n Normality) String() string {
	switch n {
	case NormallyOpen:
		return "normally open"
	case NormallyClosed:
		return "normally closed"
	default:
		return "undefined"
	}
}

type Relay struct {
	Pin        pin.Pin
	number     int
	Normality  Normality
	DeviceName string
	DeviceTypeInfo
}

func (r *Relay) Name() string {
	return r.DeviceName
}

func (r *Relay) On() {
	switch normality := r.Normality; normality {
	case NormallyOpen:
		log.Debug().Msgf("setting pin %+v to high", uint8(r.Pin.Pin))
		r.Pin.Output()
		r.Pin.High()
	case NormallyClosed:
		log.Debug().Msgf("setting pin %+v to low", uint8(r.Pin.Pin))
		r.Pin.Output()
		r.Pin.Low()
	default:
	}
}

func (r *Relay) Off() {
	switch normality := r.Normality; normality {
	case NormallyOpen:
		r.Pin.Output()
		r.Pin.Low()
	case NormallyClosed:
		r.Pin.Output()
		r.Pin.High()
	default:
	}
}

func (r *Relay) Number() int {
	return r.number
}

func (r *Relay) State() (bool, error) {
	err := rpio.Open()
	if err != nil {
		return false, fmt.Errorf("unable to open rpio")
	}
	defer rpio.Close()
	state := r.Pin.Read()
	switch r.Normality {
	case NormallyClosed:
		if state == rpio.High {
			return false, nil
		}
		return true, nil
	case NormallyOpen:
		if state == rpio.High {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("how did you get to this state")
}

func (r *Relay) TypeInfo() DeviceTypeInfo {
	return r.DeviceTypeInfo
}

func NewRelay(deviceName string, pinNumber int, relayNumber int, normality Normality, category DeviceCategory) *Relay {
	log.Debug().Msgf("creating relay for device: %s with pinNumber %d, and normality: %s deviceType: %s", deviceName, pinNumber, normality.String(), category)
	return &Relay{DeviceName: deviceName, Pin: *pin.NewPin(pinNumber), number: relayNumber, Normality: normality, DeviceTypeInfo: DeviceTypeInfo{Category: category}}
}
