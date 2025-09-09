package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// HTTP client that connects via Unix socket
func newUnixClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
		Timeout: 5 * time.Second,
	}
}

// sends a PUT request with key and value
func sendPUT(client *http.Client, key, value string) error {
	url := fmt.Sprintf("http://unix/put?key=%s&value=%s", key, value)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	socketPath := "/tmp/rackkv.sock"
	client := newUnixClient(socketPath)

	// Benchmark parameters
	totalRequests := 1000000
	concurrency := 100

	fmt.Printf("Running %d PUT requests with concurrency %d...\n", totalRequests, concurrency)

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency) // limit concurrent goroutines
	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- struct{}{} 

		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d", i)
			if err := sendPUT(client, key, value); err != nil {
				fmt.Printf("Error PUT %s: %v\n", key, err)
			}
			<-sem 
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Completed %d PUT requests in %v\n", totalRequests, duration)
	fmt.Printf("Average: %.2f ms per request\n", float64(duration.Milliseconds())/float64(totalRequests))
}
