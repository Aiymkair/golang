package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

// IsRetryable
func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// CalculateBackoff
func CalculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if backoff > maxDelay {
		backoff = maxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

// ExecutePayment
func ExecutePayment(ctx context.Context, client *http.Client, req *http.Request,
	maxRetries int, baseDelay, maxDelay time.Duration) (*http.Response, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}

		if !IsRetryable(resp, err) {
			return resp, fmt.Errorf("non-retryable error: %w", err)
		}

		if attempt == maxRetries-1 {
			return resp, fmt.Errorf("last attempt failed: %w", err)
		}

		wait := CalculateBackoff(attempt, baseDelay, maxDelay)
		fmt.Printf("Attempt %d failed, waiting %v...\n", attempt+1, wait)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}
	return nil, fmt.Errorf("max retries exceeded")
}

func main() {
	attemptCounter := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCounter++
		if attemptCounter <= 3 {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	const (
		maxRetries = 5
		baseDelay  = 500 * time.Millisecond
		maxDelay   = 5 * time.Second
	)

	fmt.Println("=== Starting payment execution ===")
	resp, err := ExecutePayment(ctx, client, req, maxRetries, baseDelay, maxDelay)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Attempt %d: Success! Response: %s\n", attemptCounter, body)
}
