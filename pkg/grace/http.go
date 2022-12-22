package grace

import (
	"context"
	"net/http"
	"sync"

	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"go.uber.org/zap"
)

type HTTPGrace struct {
	env  *baseenv.Environment
	svr  *http.Server
	addr string
}

func New(env *baseenv.Environment, mux http.Handler, addr string) *HTTPGrace {
	g := &HTTPGrace{
		env:  env,
		addr: addr,
	}
	g.svr = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return g
}

func (g HTTPGrace) Start() error {
	env := g.env

	GoWG(func() error {
		err := g.svr.ListenAndServe()
		env.Logger().Info("HTTP server closed")

		return err
	}, func(err error) {
		if err != nil && err != http.ErrServerClosed {
			env.Logger().Fatal("HTTP Server crashed, err is:" + err.Error())
		}
	}, env.CancelWG())

	env.Logger().Info("Listening HTTP", zap.String("addr", g.addr))

	return nil
}

func (g HTTPGrace) GraceStop(ctx context.Context) error {
	return g.svr.Shutdown(ctx)
}

type ErrorHandler func(error)
type Job func() error

func GoWG(job Job, catch ErrorHandler, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		catch(job())
	}()
}
