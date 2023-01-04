package apiserver

import (
	"context"
	"log"

	"github.com/ramseyjiang/go_senior_to_principle/internal/db"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/internal/server"
	"github.com/ramseyjiang/go_senior_to_principle/pkg/grace"
)

func Main(dbType string) {
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

	log.Println(graceServerList)

	err = db.InitDB(gEnv, dbType)
	if err != nil {
		log.Fatal(err)
	}

	cancel()
	gEnv.Logger().Info("Server exiting")
	gEnv.CloseWG().Wait()
	_ = gEnv.Close()
}
