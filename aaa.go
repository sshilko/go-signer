package main

import (
	"fmt"
	"sync"
)

func main() {
	Example()
	fmt.Println("all done!")
}

func MergeStreams(s1, s2 chan int) chan int {
	mychan := make(chan int)

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		for v := range s1 {
			mychan <- v
		}
		fmt.Println("s1 finished")
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for v := range s2 {
			mychan <- v
		}
		fmt.Println("s2 finished")
		defer wg.Done()
	}()

	go func(w *sync.WaitGroup, myc chan int) {
		w.Wait()
		close(myc)
	}(wg, mychan)

	return mychan
}

func Example() {
	a := make(chan int)
	b := make(chan int)
	c := MergeStreams(a, b)

	go func() {
		a <- 1
		a <- 2
		close(a)
	}()

	go func() {
		b <- 3
		b <- 4
		close(b)
	}()

	// Expected output:
	// 1 2 3 4 in any order
	for v := range c {
		fmt.Println(v)
	}

}
