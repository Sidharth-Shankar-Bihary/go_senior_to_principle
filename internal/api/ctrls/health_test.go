package ctrls

import (
	"net/http"
	"testing"
)

func TestCheckHealth(t *testing.T) {
	RunAPITests(t, []APITestCase{
		{
			tag:      "Health Check",
			method:   "GET",
			route:    "/api/health",
			url:      "/api/health",
			body:     `{"data":"Health is ok."}`,
			function: CheckHealth,
			status:   http.StatusOK,
		},
	})
}
