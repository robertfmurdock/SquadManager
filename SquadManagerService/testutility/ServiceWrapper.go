package testutility

import (
	"encoding/json"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"testing"
	"net/http"
	"io"
	"bytes"
)

type ServiceWrapper struct {
	t *testing.T
	Handler http.Handler
}

func Wrap(t *testing.T, handler http.Handler) *ServiceWrapper {
	return &ServiceWrapper{t, handler}
}

func NewRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func getPostBody(body interface{}) ([]byte, error) {
	if value, ok := body.(string); ok {
		return []byte(value), nil
	} else {
		return json.Marshal(body)
	}
}

func MakePostRequest(t *testing.T, url string, body interface{}) *http.Request {
	value, err := getPostBody(body)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest("POST", url, bytes.NewReader(value))

	if err != nil {
		t.Fatal(err)
	}
	return request
}

func (self ServiceWrapper) PerformPostSquadAndGetId() string {
	response := self.PerformPostSquad(self.t)
	assert.Equal(self.t, http.StatusAccepted, response.Code)
	var newSquadId string
	json.Unmarshal(response.Body.Bytes(), &newSquadId)
	return newSquadId
}

func (self ServiceWrapper) PerformGetSquad(squadId string) api.Squad {
	request := NewRequest(self.t, "GET", "/squad/"+squadId, nil)

	recorder := self.getResponse(request)

	assert.Equal(self.t, recorder.Code, 200)

	var loadedJson api.Squad
	err := json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	if err != nil {
		self.t.Fatal(err)
	}

	return loadedJson
}

func (self ServiceWrapper) PerformPostSquadMember(squadId string, squadMember api.SquadMember) string {
	request := MakePostRequest(self.t, "/squad/"+squadId, squadMember)

	response := self.getResponse(request)

	assert.Equal(self.t, http.StatusAccepted, response.Code)
	var newSquadMemberId string
	json.Unmarshal(response.Body.Bytes(), &newSquadMemberId)

	return newSquadMemberId
}

func (self ServiceWrapper) PerformPostSquad(t *testing.T) *httptest.ResponseRecorder {
	request := MakePostRequest(t, "/squad", "")

	return self.getResponse(request)
}

func (self ServiceWrapper) getResponse(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	self.Handler.ServeHTTP(recorder, request)
	return recorder
}

func (self ServiceWrapper) PerformGetSquadList(t *testing.T) []string {
	request := NewRequest(t, "GET", "/squad", nil)

	recorder := self.getResponse(request)

	assert.Equal(t, recorder.Code, 200)

	var loadedJson []string
	json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	return loadedJson
}