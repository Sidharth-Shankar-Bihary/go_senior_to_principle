package server

import (
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/grace"
)

func InitServer(gEnv *env.Environment) ([]grace.Closer, error) {
	gEnv.Log.Info("server start")
	handler, err := provideHTTPServer(gEnv)
	if err != nil {
		return nil, err
	}

	return provideGraceServices(gEnv, handler), nil
}
