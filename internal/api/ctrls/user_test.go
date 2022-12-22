package ctrls

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	body := `{
				"status": 200,
				"user": {
					"id":1,
					"created_at":"0001-01-01T00:00:00Z",
					"updated_at":"0001-01-01T00:00:00Z",
					"deleted_at": null,
					"first_name":"",
					"last_name":"",
					"address": "",
					"email":""
				}
			}`
	tests := []APITestCase{
		{
			"test 1 - get a User",
			"GET",
			"/users/:id",
			"/users/1",
			body,
			GetUser,
			http.StatusOK,
		},
	}
	RunAPITests(t, tests)
}
