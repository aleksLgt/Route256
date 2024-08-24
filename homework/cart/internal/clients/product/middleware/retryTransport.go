package middleware

import (
	"fmt"
	"net/http"
)

type RetryTransport struct {
	Transport  http.RoundTripper
	MaxRetries int
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	const enhancementYourCalmStatus = 420

	for i := 0; i < t.MaxRetries; i++ {
		res, err := t.Transport.RoundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("RoundTrip failed: %w", err)
		}

		if res.StatusCode != enhancementYourCalmStatus && res.StatusCode != http.StatusTooManyRequests {
			return res, nil
		}
	}

	return nil, fmt.Errorf("max retries exceeded")
}
