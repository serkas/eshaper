package main

import "github.com/serkas/eshaper"

func main() {

	rateUpdateCh := make(chan int64, 0)

	shaper := eshaper.NewShaper(2, rateUpdateCh)
	shaper.Run()
}
