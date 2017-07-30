package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/stretchr/testify/assert"
)

type ServiceWrapper struct {
	t       *testing.T
	Handler http.Handler
}

func Wrap(t *testing.T, handler http.Handler) *ServiceWrapper {
	return &ServiceWrapper{t, handler}
}

func (wrapper *ServiceWrapper) PerformRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	wrapper.Handler.ServeHTTP(recorder, request)
	return recorder
}

func (wrapper *ServiceWrapper) PerformPostSquad() string {
	var newSquadId string
	PostSquad().
		WithStatus(http.StatusAccepted).
		LoadJson(&newSquadId).
		Perform(wrapper)
	return newSquadId
}

func (wrapper *ServiceWrapper) PerformGetSquad(squadId string) api.Squad {
	squad := api.Squad{}
	GetSquad(squadId).
		WithStatus(http.StatusOK).
		LoadJson(&squad).
		Perform(wrapper)
	return squad
}

func (wrapper *ServiceWrapper) PerformGetSquadList() []string {
	var loadedJson []string
	GetSquadList().
		WithStatus(http.StatusOK).
		LoadJson(&loadedJson).
		Perform(wrapper)
	return loadedJson
}

func (wrapper *ServiceWrapper) PerformPostSquadMember(squadId string, squadMember api.SquadMember) string {
	var newSquadMemberId string
	PostSquadMember(squadId, squadMember).
		WithStatus(http.StatusAccepted).
		LoadJson(&newSquadMemberId).
		Perform(wrapper)
	return newSquadMemberId
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

type Perform func(wrapper *ServiceWrapper) *httptest.ResponseRecorder

func (perform Perform) Then(checker ResponseChecker) Perform {
	return func(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
		response := perform(wrapper)
		checker(wrapper.t, response)
		return response
	}
}

func (perform Perform) Perform(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
	return perform(wrapper)
}

func (perform Perform) WithStatus(status int) Perform {
	return perform.Then(WithStatus(status))
}

func (perform Perform) LoadJson(value interface{}) Perform {
	return perform.Then(LoadJson(value))
}

type ResponseChecker func(_ *testing.T, _ *httptest.ResponseRecorder) *httptest.ResponseRecorder

func WithStatus(status int) ResponseChecker {
	return func(t *testing.T, response *httptest.ResponseRecorder) *httptest.ResponseRecorder {
		assert.Equal(t, status, response.Code)
		return response
	}
}

func LoadJson(loadedJson interface{}) ResponseChecker {
	return func(t *testing.T, response *httptest.ResponseRecorder) *httptest.ResponseRecorder {
		err := json.Unmarshal(response.Body.Bytes(), loadedJson)
		if err != nil {
			t.Fatal(err)
		}
		return response
	}
}

func GetSquadList() Perform {
	return func(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
		request := NewRequest(wrapper.t, "GET", "/squad", nil)
		return wrapper.PerformRequest(request)
	}
}

func GetSquad(squadId string) Perform {
	return func(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
		request := NewRequest(wrapper.t, "GET", "/squad/"+squadId, nil)
		return wrapper.PerformRequest(request)
	}
}

func PostSquad() Perform {
	return func(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
		request := MakePostRequest(wrapper.t, "/squad", "")
		return wrapper.PerformRequest(request)
	}
}
func PostSquadMember(squadId string, member api.SquadMember) Perform {
	return func(wrapper *ServiceWrapper) *httptest.ResponseRecorder {
		request := MakePostRequest(wrapper.t, "/squad/"+squadId, member)
		return wrapper.PerformRequest(request)
	}
}
