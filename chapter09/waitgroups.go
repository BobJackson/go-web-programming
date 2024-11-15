package main

import (
	"fmt"
	"sync"
	"time"
)

func printNumber3(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%d ", i)
	}
	wg.Done()
}

func printLetters3(wg *sync.WaitGroup) {
	for i := 'A'; i < 'A'+10; i++ {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%c ", i)
	}
	wg.Done()
}

func printNumber4(w chan bool) {
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%d ", i)
	}
	w <- true
}

func printLetters4(w chan bool) {
	for i := 'A'; i < 'A'+10; i++ {
		time.Sleep(1 * time.Microsecond)
		fmt.Printf("%c ", i)
	}
	w <- true
}

func usingWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(2)
	go printNumber3(&wg)
	go printLetters3(&wg)
	wg.Wait()
}

func usingChannel() {
	w1, w2 := make(chan bool), make(chan bool)
	go printNumber4(w1)
	go printLetters4(w2)
	<-w1
	<-w2
}

func thrower(c chan int) {
	for i := 0; i < 5; i++ {
		c <- i
		fmt.Println("Threw  >>", i)
	}
}
func catcher(c chan int) {
	for i := 0; i < 5; i++ {
		num := <-c
		fmt.Println("Caught <<", num)
	}
}

func main() {
	//usingWaitGroup()
	//usingChannel()

	c := make(chan int)
	go thrower(c)
	go catcher(c)
	time.Sleep(100 * time.Microsecond)
}
