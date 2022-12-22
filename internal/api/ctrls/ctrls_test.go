package ctrls

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type APITestCase struct {
	tag      string
	method   string
	route    string
	url      string
	body     string
	function gin.HandlerFunc
	status   int
}

// Creates new router in testing mode
func createRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	return router
}

// CommonTestAPI is used to run single API test case. It makes HTTP request and returns its response
func newTestAPI(router *gin.Engine, method string, route string, url string, function gin.HandlerFunc, body string) *httptest.ResponseRecorder {
	// registers a new request handle with the given method, route and real function
	router.Handle(method, route, function)

	// httptest.NewRecorder() will be responsible for writing the response for a route
	resp := httptest.NewRecorder()

	// http.NewRequest returns a new Request with given a method, URL, and body.
	req, _ := http.NewRequest(method, url, bytes.NewBufferString(body))
	router.ServeHTTP(resp, req)
	return resp
}

func RunAPITests(t *testing.T, tests []APITestCase) {
	for _, test := range tests {
		router := createRouter()
		resp := newTestAPI(router, test.method, test.route, test.url, test.function, test.body)
		assert.Equal(t, test.status, resp.Code)
		assert.JSONEq(t, test.body, resp.Body.String())
	}
}
