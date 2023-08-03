package utils

import (
	"math/rand"
	"time"
)

type Spread struct {
	Min    int
	Max    int
	source rand.Source
}

func NewSpread(min int, max int) *Spread {
	src := rand.NewSource(time.Now().UnixNano())
	return &Spread{
		Min:    min,
		Max:    max,
		source: src,
	}
}

func (s *Spread) IsValid() bool {
	return s.Min < s.Max
}

func (s *Spread) Rand() int {
	return rand.New(s.source).Intn(s.Max-s.Min) + s.Min
}
