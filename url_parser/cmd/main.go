package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var urls = []string{
		"https://hackernoon.com",
		"https://www.facebook.com",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://www.netflix.com/pl-en",
		"http://cccc",
		"https://www.booking.com",
	}

	resChan := make(chan string, len(urls))

	wg := &sync.WaitGroup{}
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			_, err := http.Get(url)

			if err != nil {
				resChan <- fmt.Sprintf("Url not ok: %s", url)
				return
			}

			resChan <- fmt.Sprintf("Url: %s", url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for u := range resChan {
		fmt.Println(u)
	}
}
