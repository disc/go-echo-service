package main

import (
	"context"
	"flag"
	"github.com/disc/go-echo-service/pkg/echoservice"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	listen = flag.String("listen", ":8080", "HTTP listen address")
	proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy echo requests")
)

func main() {
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

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

	var svc echoservice.EchoService
	svc = echoservice.NewEchoService()
	svc = echoservice.ProxyingMiddleware(context.Background(), *proxy, logger)(svc)
	svc = echoservice.LoggingMiddleware(logger)(svc)
	svc = echoservice.InstrumentingMiddleware(requestCount, requestLatency)(svc)

	echoHandler := httptransport.NewServer(
		echoservice.MakeEchoEndpoint(svc),
		echoservice.DecodeEchoRequest,
		echoservice.EncodeResponse,
	)

	router := http.NewServeMux()
	router.Handle("/echo", echoHandler)
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    *listen,
		Handler: router,
	}

	go func() {
		logger.Log("msg", "HTTP", "addr", *listen)
		logger.Log("err", server.ListenAndServe())
	}()

	graceful(server, logger, 5*time.Second)
}

func graceful(server *http.Server, logger log.Logger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Log("msg", "Shutdown with timeout", "timeout", timeout)

	if err := server.Shutdown(ctx); err != nil {
		logger.Log("error", err)
	} else {
		logger.Log("msg", "Server gracefully stopped")
	}
}
