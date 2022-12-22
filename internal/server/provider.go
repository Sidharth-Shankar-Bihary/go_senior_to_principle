package server

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/internal/middleware"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/grace"
)

func provideGraceServices(e *env.Environment, httpSrv http.Handler) []grace.Closer {
	httpPort := e.C().Meta.HTTPPort
	if httpPort == 0 {
		httpPort = 8888
	}

	return []grace.Closer{grace.New(e.Environment, httpSrv, ":"+strconv.Itoa(httpPort))}
}

func provideHTTPServer(e *env.Environment) (http.Handler, error) {
	hMux := provideHMux(e)

	// convert router to http.Handler, which is used in middleware.InjectTracer
	router := api.NewRouter().Handler()

	// make all routes can be traced using jaeger
	hMux.Handle("/", middleware.InjectTracer(e)(router))

	return hMux, nil
}

func provideHMux(e *env.Environment) *http.ServeMux {
	hMux := http.NewServeMux()

	hMux.HandleFunc("/live", func(writer http.ResponseWriter, request *http.Request) {
		e.Logger().Debug("live")
		_, _ = writer.Write([]byte("api server mux live"))
	})
	hMux.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		e.Logger().Debug("ready")
		_, _ = writer.Write([]byte("api server mux ready"))
	})
	// *****************************************************
	// pprof
	hMux.Handle("private/debug/", http.DefaultServeMux)
	// *****************************************************

	// prom http
	hMux.Handle("/metrics", promhttp.Handler())

	return hMux
}
