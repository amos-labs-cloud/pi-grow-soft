package pin

import (
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

type Pin struct {
	rpio.Pin
}

type Service struct {
}

func New(pins ...uint8) *Service {
	return &Service{}
}

func NewPin(pin uint8) *Pin {
	return &Pin{rpio.Pin(pin)}
}

func (p *Pin) High() {
	log.Debug().Msgf("Setting %+v to high", p.Pin)
	p.Pin.Output()
	p.Pin.High()

}

func (p *Pin) Low() {
	log.Debug().Msgf("Setting %+v to low", p.Pin)
	p.Pin.Output()
	p.Pin.Low()
}
