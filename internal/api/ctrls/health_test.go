package ctrls

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckHealth(t *testing.T) {
	tests := []APITestCase{
		{
			tag:      "Health Check",
			method:   "GET",
			route:    "/health",
			url:      "/health",
			body:     `{"data":"Health is ok."}`,
			function: h.CheckHealth,
			status:   http.StatusOK,
		},
	}

	for _, test := range tests {
		resp := newTestAPI(router, test.method, test.route, test.url, test.function, test.body)
		assert.Equal(t, test.status, resp.Code)
		assert.JSONEq(t, test.body, resp.Body.String())
	}
}
