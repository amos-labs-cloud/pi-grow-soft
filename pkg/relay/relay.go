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
	number     uint8
	Normality  Normality
	DeviceName string
}

func (r *Relay) Name() string {
	return r.DeviceName
}

func (r *Relay) On() {
	switch normality := r.Normality; normality {
	case NormallyOpen:
		log.Debug().Msgf("setting pin %r to high", uint8(r.Pin.Pin))
		r.Pin.Output()
		r.Pin.High()
	case NormallyClosed:
		log.Debug().Msgf("setting pin %r to low", uint8(r.Pin.Pin))
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

func (r *Relay) Number() uint8 {
	return r.number
}

func (r *Relay) State() (bool, error) {
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

func NewRelay(deviceName string, pinNumber uint8, relayNumber uint8, normality Normality) *Relay {
	log.Debug().Msgf("creating relay for device: %s with pinNumber %d, and normality: %s", deviceName, pinNumber, normality.String())
	return &Relay{DeviceName: deviceName, Pin: *pin.NewPin(pinNumber), number: relayNumber, Normality: normality}
}
