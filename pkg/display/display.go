package display

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type LCDType uint

const (
	LCD1602 LCDType = iota
)

type Display interface {
	WriteLines(lines ...string)
}

type DefaultViewFunc interface {
	LoadDefaultPage() func() []*string
}

type Page []string

type Service struct {
	display Display
	pages   []PageFunc
}

func New(display Display) *Service {
	return &Service{
		display: display,
	}
}

func (s *Service) Clear() {
	s.display.WriteLines("", "")
}

type PageFunc func() Page

func (s *Service) ClearPages() {
	s.pages = []PageFunc{}
}

func debugLines(lines ...string) {
	for _, line := range lines {
		log.Debug().Msgf("Length of this line: %s is %d", line, len(line))
	}
}

func (s *Service) AddPage(lines ...string) {
	debugLines(lines...)
	s.pages = append(s.pages, func() Page { return lines })
}

func (s *Service) AddPageFunc(pageFunc PageFunc) {
	s.pages = append(s.pages, pageFunc)
}

func (s *Service) RunPages(ctx context.Context) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			for _, page := range s.pages {
				s.Clear()
				s.display.WriteLines(page()...)
				time.Sleep(5 * time.Second)
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return cancel
}

func (s *Service) WriteDisplay(lines ...string) {
	debugLines(lines...)
	s.display.WriteLines(lines...)
}
