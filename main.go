package main

import (
	"encoding/json"
	"net/http"
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/lovababu/go-coes-poc/service"
	"github.com/lovababu/go-coes-poc/api"
	"os"
	"flag"
)

func main() {

	logger := kitlog.NewLogfmtLogger(os.Stderr)

	var svc service.StringService
	svc = service.New()

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "avol",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "avol",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "avol",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	svc = api.LoggingMiddleware{
		Logger: logger,
		Next: svc,

	}

	svc = api.InstrumentingMiddleware{
		RequestCount: requestCount,
		RequestLatency: requestLatency,
		CountResult: countResult,
		Next: svc,
	}



	uppercaseHandler := httptransport.NewServer(
		api.MakeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		api.MakeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	flag.Parse()
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(*addr, nil))
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
