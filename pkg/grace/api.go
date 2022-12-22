package grace

import (
	"context"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	var ctxClosing context.Context
	select {
	case sig := <-quit:
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
		err := service.GraceStop(ctxClosing)
		logger.Fatal("failed to stop, err is:" + err.Error())
	}
}
