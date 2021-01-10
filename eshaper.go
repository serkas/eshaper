package eshaper

// shaper generates some token messages in a special channel
// workers (consumers of the real queue) need to take a token from the shaper channel to consume one message
// so workers can not consume more payload messages than there are generated tokens

import (
	"math"
	"time"
)

type tickSettings struct {
	tickInterval  time.Duration
	tokensPerTick int64
}

type Shaper interface {
	SetRate(rate int64, interval time.Duration)
	Pass()
}

type shaper struct {
	tokenCh       chan struct{}
	tickInterval  time.Duration
	tokensPerTick int64
	rateCh        chan tickSettings
}

func New(rate int64, baseInterval time.Duration) *shaper {
	count, interval := selectInterval(rate, baseInterval)
	sh := &shaper{
		tokenCh:       make(chan struct{}, 1),
		tickInterval:  interval,
		tokensPerTick: count,
		rateCh:        make(chan tickSettings),
	}

	go sh.run()

	return sh
}

func (s *shaper) SetRate(rate int64, baseInterval time.Duration) {
	count, tInt := selectInterval(rate, baseInterval)
	s.rateCh <- tickSettings{
		tokensPerTick: count,
		tickInterval:  tInt,
	}
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
			for i := int64(0); i < s.tokensPerTick; i++ {
				s.tokenCh <- struct{}{}
			}
		}
	}
}

func (s *shaper) Pass() {
	<-s.tokenCh
}

func selectInterval(rate int64, interval time.Duration) (int64, time.Duration) {
	rps := float64(rate) / interval.Seconds()

	// On high rates, add more then one token per tick
	if rps > 2000 {
		tick := 10 * time.Millisecond
		count := int64(math.Ceil(rps * time.Second.Seconds() / tick.Seconds()))
		return count, tick
	}

	// If the rate is not so high and the provided interval is exact (one per `interval`) use the interval
	if rate == 1 {
		return rate, interval
	}

	return 1, interval / time.Duration(rate)
}
