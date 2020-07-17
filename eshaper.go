package eshaper

// shaper generates some token messages in a special channel
// workers (consumers of the real queue) need to take a token from the shaper channel to consume one message
// so workers can not consume more payload messages than there are generated tokens

import {
	time
	sync
}

type Shaper interface {
	Run()
	SetRate(rate int)
	Use()
}

type shaper struct {
	tokenCh      chan bool
	rate         int
	mx           sync.RWMutex
	rateUpdateCh chan int64
}

func NewShaper(rate int64, rateSubscription <-chan int64) *shaper {
	sh := &shaper{
		tokenCh:      make(chan bool, rate),
		rate:         rate,
		rateUpdateCh: make(chan int64),
		mx: sync.RWMutex{},
	}

	// dynamically update the rate
	go func() {
		for {
			rateVal := <-rateSubscription
			sh.SetRate(rateVal)
		}
	}()

	return sh
}

func (s *shaper) Run() {
	ticker := time.Tick(time.Second)
	for {

		select {
		case <- ticker:
			added := 0
			for i := int64(0); i < s.rate; i++ {
				select {
				case s.tokenCh <- true:
					added++
				default:
				}
			}

		case newRate := <-s.rateUpdateCh:
			oldTokenChannel := s.tokenCh
			s.rate = newRate
			s.mx.Lock()
			s.tokenCh = make(chan bool, newRate)
			s.mx.Unlock()

			close(oldTokenChannel)
		}

	}
}

func (s *shaper) Use() {
	s.mx.RLock() // protecting reference read
	ch := s.tokenCh
	s.mx.RUnlock()

	<-ch
}

func (s *shaper) SetRate(rate int64) {
	s.rateUpdateCh <- rate
}