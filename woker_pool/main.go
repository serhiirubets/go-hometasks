package main

import (
	"fmt"
	"sync"
	"time"
)

func runNumSquareWorker(jobs <-chan int, results chan<- int) {
	for j := range jobs {
		time.Sleep(1 * time.Second)
		results <- j * j
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	workerPoolNum := 3

	jobs := make(chan int, len(nums))
	results := make(chan int, len(nums))

	// WaitGroup для отслеживания завершения воркеров
	var wg sync.WaitGroup

	for w := 1; w <= workerPoolNum; w++ {
		wg.Add(1)
		go func() {
			runNumSquareWorker(jobs, results)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()      // Ждем, пока все воркеры завершат работу
		close(results) // Закрываем канал результатов
	}()

	for _, n := range nums {
		jobs <- n
	}

	close(jobs)

	for range nums {
		r := <-results
		fmt.Println(r)
	}
}
