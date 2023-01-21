package grace

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ramseyjiang/go_senior_to_principle/pkg/baseenv"
	"go.uber.org/zap"
)

type Closer interface {
	Start() error
	GraceStop(context.Context) error
}

func Serve(env *baseenv.Environment, services ...Closer) {
	ctx := env.Context()
	logger := env.Logger()

	for _, service := range services {
		err := service.Start()
		if err != nil {
			logger.Fatal("failed to start, err is:" + err.Error())
		}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	var ctxClosing context.Context
	// if necessary, hot reload can be added at here.
	select {
	case sig := <-stopChan:
		logger.Info("Signal received, shutting down server...", zap.String("signal", sig.String()))

		cc, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		ctxClosing = cc
	case <-ctx.Done():
		logger.Info("Context cancelled, shutting down server...")

		cc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ctxClosing = cc
	}

	for _, service := range services {
		_ = service.GraceStop(ctxClosing)
	}
}

type ContextJob func(context.Context) error

func ContextJobGrace(fn ContextJob) Closer {
	return &contextGrace{fn: fn}
}

type contextGrace struct {
	fn      ContextJob
	cancel  context.CancelFunc
	errChan <-chan error
}

func (c *contextGrace) Start() error {
	if c.cancel != nil {
		return errors.New("can not start twice")
	}

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	c.errChan = GoChan(WithContext(ctx, c.fn))
	return nil
}

func (c *contextGrace) GraceStop(ctx context.Context) error {
	c.cancel()
	return <-c.errChan
}

type Job func() error

func GoChan(job Job) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- job()
	}()

	return ch
}

func WithContext(ctx context.Context, cj ContextJob) Job {
	return func() error {
		return cj(ctx)
	}
}
