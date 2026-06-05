package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TotalRequests = 20000
	Concurrency   = 100
	Endpoint      = "http://localhost:4566/_cloudstack/health" // Super fast endpoint
)

func main() {
	fmt.Printf("Starting Load Test: %d requests across %d concurrent workers...\n", TotalRequests, Concurrency)
	fmt.Printf("Target: %s\n", Endpoint)

	// Create an optimized HTTP client
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = Concurrency
	t.MaxIdleConnsPerHost = Concurrency
	t.IdleConnTimeout = 30 * time.Second

	client := &http.Client{
		Transport: t,
		Timeout:   5 * time.Second,
	}

	var successCount int32
	var errorCount int32

	reqsPerWorker := TotalRequests / Concurrency

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(Concurrency)

	for i := 0; i < Concurrency; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < reqsPerWorker; j++ {
				resp, err := client.Get(Endpoint)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					continue
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				if resp.StatusCode == 200 {
					atomic.AddInt32(&successCount, 1)
				} else {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	rps := float64(TotalRequests) / duration.Seconds()
	avgLatency := (duration.Seconds() * 1000) / float64(TotalRequests) * float64(Concurrency) // Approx per req

	fmt.Printf("\n--- Results ---\n")
	fmt.Printf("Time taken:    %.2f seconds\n", duration.Seconds())
	fmt.Printf("Total reqs:    %d\n", TotalRequests)
	fmt.Printf("Successful:    %d\n", successCount)
	fmt.Printf("Failed:        %d\n", errorCount)
	fmt.Printf("Throughput:    %.2f requests/sec\n", rps)
	fmt.Printf("Avg Latency:   %.2f ms\n", avgLatency)
}
