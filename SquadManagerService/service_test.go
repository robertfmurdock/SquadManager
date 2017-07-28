package service

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"io"
)

func TestNoResponseOnMainUrl(t *testing.T) {
	request := newRequest(t, "GET", "/", nil)
	recorder := httptest.NewRecorder()

	http.HandlerFunc(MainHandler).ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 404)
}

func newRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}
	return request;
}

func TestCallGetSquad(t *testing.T) {
	request := newRequest(t, "GET", "/squad", nil)
	recorder := httptest.NewRecorder()

	http.HandlerFunc(MainHandler).ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 200)

	var loadedJson []string
	json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	expectedJson := []string{}

	assert.Equal(t, expectedJson, loadedJson)
}

