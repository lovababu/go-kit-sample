package api

import (
	"github.com/go-kit/kit/log"
	"time"
	"github.com/lovababu/go-coes-poc/service"
)


type LoggingMiddleware struct {
	Logger log.Logger
	Next   service.StringService
}

func (m LoggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		m.Logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = m.Next.Uppercase(s)
	return
}

func (m LoggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		m.Logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = m.Next.Count(s)
	return
}
