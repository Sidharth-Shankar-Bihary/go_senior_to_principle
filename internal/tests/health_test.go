package tests

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	testURL := "http://localhost:8888/health"
	resp, err := http.Get(testURL)
	if err != nil {
		log.Fatalln(err)
	}

	respBodyExpected := `{"data":"Health is ok."}`

	byteBody, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, respBodyExpected, string(byteBody))
}
