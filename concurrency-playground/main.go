package main

import (
	"fmt"
	"math/rand"
	"time"
)

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)

	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

func boring(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}

func f(left chan<- int, right <-chan int) {
  left <- 1 + <- right
}

func main() {
	// c := fanIn(boring("Joe"), boring("Ann"))
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(<-c)
	// }
	// fmt.Println("You're both boring; I'm leaving.")

  const n = 10000
  leftmost := make(chan int)
  right := leftmost
  left := leftmost
  for i := 0; i < n; i++ {
    right = make(chan int)
    go f(left, right)
    left = right
  }

  // time.Sleep(time.Second)

  go func(c chan int) { c <- 1 }(right)
  fmt.Println(<-leftmost)
}
