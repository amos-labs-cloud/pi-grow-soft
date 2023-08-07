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

func New(pins ...int) *Service {
	return &Service{}
}

func NewPin(pin int) *Pin {
	return &Pin{rpio.Pin(pin)}
}

func (p *Pin) High() {
	log.Debug().Msgf("Setting %+v to high", p.Pin)
	err := rpio.Open()
	if err != nil {
		log.Debug().Msgf("wonder how this happened??? %s", err)
	}
	defer func() {
		err := rpio.Close()
		if err != nil {
			log.Debug().Msgf("ran into this err: %s closing rpio")
		}
	}()
	p.Pin.Output()
	p.Pin.High()

}

func (p *Pin) Low() {
	log.Debug().Msgf("Setting %+v to low", p.Pin)
	err := rpio.Open()
	if err != nil {
		log.Debug().Msgf("wonder how this happened??? %s", err)
	}
	defer func() {
		err := rpio.Close()
		if err != nil {
			log.Debug().Msgf("ran into this err: %s closing rpio")
		}
	}()
	p.Pin.Output()
	p.Pin.Low()
}
