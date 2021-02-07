package main

import (
	"fmt"
	"time"

	"github.com/serkas/shaper"
)

func main() {
	events := make(chan int)

	// A source of event where you are not able to set a fixed rate. In real application it can be a queue consumer
	go func() {
		for i := 0; i < 10; i ++ {
			events <- i
		}
		close(events)
	}()
	//

	sh := shaper.New(1, time.Second)

	for e := range events {
		sh.Pass()
		fmt.Printf("%s Event %d \n", time.Now().Format(time.StampMilli), e)
	}
}
