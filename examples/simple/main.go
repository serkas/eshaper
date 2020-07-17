package main

import "../eshaper"

func main() {

	rateUpdateCh := make(chan int64, 0)

	shaper := eshaper.NewShaper(2, rateUpdateCh)
}
