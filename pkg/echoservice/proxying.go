package echoservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"net/url"
	"strings"
	"time"
)

type proxymw struct {
	ctx  context.Context
	next EchoService
	echo endpoint.Endpoint
}

func (mw proxymw) Echo(s string) (string, error) {
	response, err := mw.echo(mw.ctx, echoRequest{s})
	if err != nil {
		return "", err
	}
	resp := response.(echoResponse)
	if resp.Err != "" {
		return resp.V, errors.New(resp.Err)
	}
	return resp.V, nil
}

func ProxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
	// If instances is empty, don't proxy.
	if instances == "" {
		logger.Log("proxy_to", "none")
		return func(next EchoService) EchoService { return next }
	}

	var (
		qps         = 100
		maxAttempts = 3
		maxTime     = 250 * time.Millisecond

		instanceList = split(instances)
		endpointer   sd.FixedEndpointer
	)
	logger.Log("proxy_to", fmt.Sprint(instanceList))
	for _, instance := range instanceList {
		var e endpoint.Endpoint
		e = makeEchoProxy(ctx, instance)
		e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
		e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(e)
		endpointer = append(endpointer, e)
	}

	balancer := lb.NewRoundRobin(endpointer)
	balancedEndpoint := lb.Retry(maxAttempts, maxTime, balancer)

	return func(next EchoService) EchoService {
		return proxymw{ctx, next, balancedEndpoint}
	}
}

func makeEchoProxy(ctx context.Context, instance string) endpoint.Endpoint {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/echo"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeEchoResponse,
	).Endpoint()
}

func split(s string) []string {
	a := strings.Split(s, ",")
	for i := range a {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}
