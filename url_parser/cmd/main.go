package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
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
	errChan := make(chan string, len(urls))

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			resp, err := client.Get(url)

			if err != nil {
				errChan <- fmt.Sprintf("Url not ok: %s", url)
				return
			}

			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					errChan <- err.Error()
				}
			}(resp.Body)

			if err != nil {
				errChan <- fmt.Sprintf("Url not ok: %s", url)
				return
			}

			resChan <- fmt.Sprintf("Url: %s", url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	for u := range resChan {
		fmt.Println(u)
	}

	for err := range errChan {
		fmt.Println(err)
	}
}
