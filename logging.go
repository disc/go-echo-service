package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next EchoService) EchoService {
		return logwm{logger, next}
	}
}

type logwm struct {
	logger log.Logger
	next   EchoService
}

func (mw logwm) Echo(s string) (output string, err error) {
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
