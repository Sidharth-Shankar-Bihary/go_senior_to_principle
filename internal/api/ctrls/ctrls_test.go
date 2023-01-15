package ctrls

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type APITestCase struct {
	tag      string
	method   string
	route    string
	url      string
	body     string
	function func(c *gin.Context)
	status   int
}

var h *Handler
var router = createRouter()

// Creates new router in testing mode
func createRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // for example: []string{"https://foo.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	}))

	return r
}

// CommonTestAPI is used to run single API test case. It makes HTTP request and returns its response
func newTestAPI(router *gin.Engine, method string, route string, url string, function func(c *gin.Context), body string) *httptest.ResponseRecorder {
	// registers a new request handle with the given method, route and real function
	router.Handle(method, route, function)

	// httptest.NewRecorder() will be responsible for writing the response for a route
	resp := httptest.NewRecorder()

	// http.NewRequest returns a new Request with given a method, URL, and body.
	req, _ := http.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	return resp
}
