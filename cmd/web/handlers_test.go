package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.flaviogalon.github.io/internal/assert"
)

func TestPing(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	app := application{}

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.ping(responseRecorder, r)

	result := responseRecorder.Result()

	assert.Equal(t, result.StatusCode, http.StatusOK)

	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
