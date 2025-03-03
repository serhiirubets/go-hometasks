package main

import (
	"fmt"
	"sync"
)

func joinChannels(chs ...<-chan int) <-chan int {
	resChan := make(chan int, 20)

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(chs))
		for _, ch := range chs {
			go func(ch <-chan int, wg *sync.WaitGroup) {
				defer wg.Done()

				for n := range ch {
					resChan <- n
				}
			}(ch, &wg)
		}

		wg.Wait()

		close(resChan)
	}()

	return resChan
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() { a <- 1; a <- 2; close(a) }()
	go func() { b <- 3; close(b) }()
	go func() { c <- 4; c <- 5; c <- 6; close(c) }()

	for num := range joinChannels(a, b, c) {
		fmt.Println(num)
	}
}
