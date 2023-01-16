package tests

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	testURL := "http://localhost:8888/health"
	resp, err := http.Get(testURL)
	if err != nil {
		log.Fatalln(err)
	}

	respBodyExpected := `{"data":"Health is ok."}`

	byteBody, err := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, respBodyExpected, string(byteBody))
}
