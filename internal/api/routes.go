package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/ctrls"
)

type Router struct {
	router *gin.Engine
}

func (r *Router) Handler() http.Handler {
	return r.router.Handler()
}

func path(endpoint string) string {
	return fmt.Sprintf("/api/%s", endpoint)
}

func NewRouter() *Router {
	r := gin.Default()

	// use cors middleware. cors.Default() means allow all origins. If replace cors.Default() with cors.New(), you define allow by yourself.
	r.Use(cors.Default())

	// r.Handle and r.GET, these are two ways for gin routes.
	r.Handle("GET", path("health"), ctrls.CheckHealth)

	r.GET(path("users/:id"), ctrls.GetUser)

	return &Router{
		r,
	}
}
