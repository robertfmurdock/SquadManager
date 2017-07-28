package service

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestCallGetSquad(t *testing.T) {

	request, err := http.NewRequest("GET", "/squad", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	http.HandlerFunc(MainHandler).ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 200)

	var loadedJson []string
	json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	expectedJson := []string{}

	assert.Equal(t, expectedJson, loadedJson)
}

