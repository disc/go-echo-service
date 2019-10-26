package main

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	"time"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
) ServiceMiddleware {
	return func(next EchoService) EchoService {
		return instrmw{requestCount, requestLatency, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           EchoService
}

func (mw instrmw) Echo(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "echo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Echo(s)
	return

}
