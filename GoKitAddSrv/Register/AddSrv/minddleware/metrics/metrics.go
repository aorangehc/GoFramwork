package metrics

import (
	"context"
	"fmt"
	"gokitaddsrv/service"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           service.AddService
}

// NewInstrumentingMiddleware returns a service middleware that instruments
// the service with metrics.
func NewInstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram, result metrics.Histogram) func(service.AddService) service.AddService {
	return func(next service.AddService) service.AddService {
		return instrumentingMiddleware{
			requestCount:   counter,
			requestLatency: latency,
			countResult:    result,
			next:           next,
		}
	}
}

func (mw instrumentingMiddleware) Add(ctx context.Context, a, b int) (v int, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Add", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(v))
	}(time.Now())

	v, err = mw.next.Add(ctx, a, b)
	return
}

func (mw instrumentingMiddleware) Concat(ctx context.Context, a, b string) (v string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Concat", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(len(v)))
	}(time.Now())

	v, err = mw.next.Concat(ctx, a, b)
	return
}
