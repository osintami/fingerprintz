// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ThrottledTransport struct {
	roundTripper http.RoundTripper
	rateLimiter  *rate.Limiter
}

func (x *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := x.rateLimiter.Wait(r.Context()); err != nil {
		return nil, err
	}
	return x.roundTripper.RoundTrip(r)
}

func NewThrottledTransport(limitPeriod time.Duration, requestCount int, transportWrap http.RoundTripper) http.RoundTripper {
	return &ThrottledTransport{
		roundTripper: transportWrap,
		rateLimiter:  rate.NewLimiter(rate.Every(limitPeriod), requestCount),
	}
}
