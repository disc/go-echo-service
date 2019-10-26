package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"time"
)

type loggingMiddleware struct {
	logger log.Logger
	next   EchoService
}

func (mw loggingMiddleware) Echo(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "echo",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Echo(s)
	return
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           EchoService
}

func (mw instrumentingMiddleware) Echo(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "echo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.next.Echo(s)
	return

}

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "echo_service",
		Name:      "request_count",
		Help:      "Number of received requests.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "echo_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc EchoService
	svc = echoService{}
	svc = loggingMiddleware{logger, svc}
	svc = instrumentingMiddleware{requestCount, requestLatency, svc}

	echoHandler := httptransport.NewServer(
		makeEchoEndpoint(svc),
		decodeEchoRequest,
		encodeResponse,
	)

	http.Handle("/echo", echoHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP server", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
