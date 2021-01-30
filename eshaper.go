// Shaper can enforce limit on maximum number of code executions per instance of time.

package eshaper

import (
	"math"
	"time"
)

type tickSettings struct {
	tickInterval  time.Duration
	tokensPerTick int
}

type Shaper interface {
	SetMaxRate(rate int, interval time.Duration)
	Pass()
}

type shaper struct {
	tokenCh       chan struct{}
	tickInterval  time.Duration
	tokensPerTick int
	rateCh        chan tickSettings
}

// Creates a new instance of shaper.
// Parameters define the  maximum `number` of events per time `interval`.
func New(number int, interval time.Duration) *shaper {
	if number < 1 {
		number = 1
	}

	count, tickInterval := selectInterval(number, interval)
	s := &shaper{
		tokenCh:       make(chan struct{}, 1),
		tickInterval:  tickInterval,
		tokensPerTick: count,
		rateCh:        make(chan tickSettings),
	}

	go s.run()

	return s
}

// SetRate defines the maximum possible rate of calling Pass() without blocking.
// Parameters define the  maximum `number` of events per time `interval`.
func (s *shaper) SetMaxRate(number int, interval time.Duration) {
	if number < 1 {
		number = 1
	}

	count, tInt := selectInterval(number, interval)
	s.rateCh <- tickSettings{
		tokensPerTick: count,
		tickInterval:  tInt,
	}
}

// Pass limits the rate of an execution loop where it is inserted.
// If the set rate is not reached, it returns immediately. Otherwise, it blocks for some time to adjust the rates.
func (s *shaper) Pass() {
	<-s.tokenCh
}

func (s *shaper) run() {
	ticker := time.NewTicker(s.tickInterval)
	for {
		select {
		case settings := <-s.rateCh:
			s.tickInterval = settings.tickInterval
			s.tokensPerTick = settings.tokensPerTick
			ticker.Stop()
			ticker = time.NewTicker(s.tickInterval)

		case <-ticker.C:
			for i := 0; i < s.tokensPerTick; i++ {
				s.tokenCh <- struct{}{}
			}
		}
	}
}

func selectInterval(number int, interval time.Duration) (int, time.Duration) {
	rps := float64(number) / interval.Seconds()

	// On high rates, add more then one token per tick
	if rps > 2000 {
		tick := 10 * time.Millisecond
		count := int(math.Ceil(rps * time.Second.Seconds() / tick.Seconds()))
		return count, tick
	}

	// If the rate is not so high and the provided interval is exact (one per `interval`) use the interval
	if number > 1 {
		return 1, interval / time.Duration(number)
	}

	return number, interval
}
