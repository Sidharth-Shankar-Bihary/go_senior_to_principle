package main

import (
	"context"
	"log"

	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/internal/server"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/grace"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	gEnv, err := env.InitEnv(ctx, cancel)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	graceServerList, err := server.InitServer(gEnv)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	grace.Serve(gEnv.Environment, graceServerList...)
}
