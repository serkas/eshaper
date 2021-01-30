package eshaper

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShaper(t *testing.T) {
	var rps = 150
	s := New(rps, time.Second)

	var count = 100
	start := time.Now()
	for i := 0; i < count; i++ {
		s.Pass()
	}

	elapsed := time.Since(start)
	expected := time.Duration(count) * time.Second/time.Duration(rps)
	relativeDiff := (elapsed.Seconds() - expected.Seconds()) / expected.Seconds()
	var tolerance = 0.05
	assert.LessOrEqual(t, math.Abs(relativeDiff), tolerance)
}

func TestShaper_RateChange(t *testing.T) {
	var rps = 100
	s := New(rps, time.Second)

	var count = 50
	start := time.Now()
	for i := 0; i < count; i++ {
		s.Pass()
	}

	elapsed := time.Since(start)
	expected := time.Duration(count) * time.Second/time.Duration(rps)
	relativeDiff := (elapsed.Seconds() - expected.Seconds()) / expected.Seconds()
	var tolerance = 0.05
	assert.LessOrEqual(t, math.Abs(relativeDiff), tolerance)

	var rps2 = 500
	s.SetMaxRate(rps2, time.Second)

	start = time.Now()
	for i := 0; i < int(count); i++ {
		s.Pass()
	}
	elapsed = time.Since(start)
	expected = time.Duration(count) * time.Second/time.Duration(rps2)
	relativeDiff = (elapsed.Seconds() - expected.Seconds()) / expected.Seconds()
	assert.LessOrEqual(t, math.Abs(relativeDiff), tolerance)
}