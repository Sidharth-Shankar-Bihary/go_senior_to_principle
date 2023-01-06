package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ramseyjiang/go_senior_to_principle/internal/api/ctrls"
	"github.com/ramseyjiang/go_senior_to_principle/internal/env"
	"github.com/ramseyjiang/go_senior_to_principle/internal/middleware"
)

type Router struct {
	router *gin.Engine
}

func (r *Router) Handler() http.Handler {
	return r.router.Handler()
}

func NewRouter(gEnv *env.Environment, h *ctrls.Handler) *Router {
	r := gin.Default()

	// using this way to set gEnv into the gin.Context, then it will be used in every route gin.Context.
	// But it will lead hard to do test, so I don't use this way, I use the pattern dependency injection to fix it.
	// r.Use(func(c *gin.Context) {
	// 	c.Set("env", gEnv)
	// })

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// use cors middleware. cors.Default() means allow all origins.
	// Here customise a simple CORS rule.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // for example: []string{"https://foo.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	}))

	// r.Handle and r.GET, these are two ways for gin routes.
	r.Handle("GET", "health", h.CheckHealth)

	pubRoutes := r.Group("/api")
	pubRoutes.POST("register", h.Register)
	pubRoutes.POST("login", h.Login)
	pubRoutes.POST("token/refresh", h.RefreshToken)

	priRoutes := r.Group("/api/user/")
	priRoutes.Use(middleware.Auth(gEnv))
	priRoutes.POST("logout", h.Logout)
	priRoutes.GET(":id", h.GetUser)
	priRoutes.GET("current", h.CurrentUser)

	return &Router{
		r,
	}
}
