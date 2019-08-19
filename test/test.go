package main

import (
	"fmt"
	"math/rand"
	"time"
)

func wait(c chan int)  {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(5)
	time.Sleep(time.Duration(i) * time.Second)
	c <- i
}

func main()  {
	chan1 := make(chan int)
	chan2 := make(chan int)

	go wait(chan1)
	go wait(chan2)

	select {
	case i := <-chan1:
		fmt.Printf("Recived %d on chan 1", i )
	case i := <- chan2:
		fmt.Printf("Recived %d on chan 1", i )
	}
}