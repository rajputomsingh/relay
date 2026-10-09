package delivery

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryDelay(t *testing.T) {
	worker := NewWorker(nil, Config{
		MaxAttempts:    8,
		RequestTimeout: 10 * time.Second,
		LeaseDuration:  30 * time.Second,
	})

	tests := []struct {
		name    string
		attempt int
		min     time.Duration
		max     time.Duration
	}{
		{"first attempt", 1, 2 * time.Second, 2500 * time.Millisecond},
		{"second attempt", 2, 4 * time.Second, 5 * time.Second},
		{"third attempt", 3, 8 * time.Second, 10 * time.Second},
		{"capped backoff", 10, 5 * time.Minute, 6*time.Minute + 15*time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delay := worker.retryDelay(tt.attempt)

			if delay < tt.min || delay > tt.max {
				t.Fatalf(
					"retryDelay(%d) = %v; want between %v and %v",
					tt.attempt,
					delay,
					tt.min,
					tt.max,
				)
			}
		})
	}
}

func TestWorkerDoesNotFollowRedirects(t *testing.T) {
	worker := NewWorker(nil, Config{
		WebhookURL:     "https://example.com/webhook",
		RequestTimeout: 10 * time.Second,
	})

	req, err := http.NewRequest(
		http.MethodGet,
		"https://example.com/redirect",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	err = worker.client.CheckRedirect(req, nil)
	if err != http.ErrUseLastResponse {
		t.Fatalf("CheckRedirect() = %v; want %v", err, http.ErrUseLastResponse)
	}
}
