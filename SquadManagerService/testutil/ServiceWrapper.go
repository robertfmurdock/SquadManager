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

type Tester struct {
	t       *testing.T
	Handler http.Handler
}

func New(t *testing.T, handler http.Handler) *Tester {
	return &Tester{t, handler}
}

func (tester *Tester) PerformRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	tester.Handler.ServeHTTP(recorder, request)
	return recorder
}

func (tester *Tester) DoRequest(method, urlStr string, body interface{}) Response {
	value, err := getPostBody(body)
	if err != nil {
		tester.t.Fatal(err)
	}
	bodyReader := bytes.NewReader(value)
	request := newRequest(tester.t, method, urlStr, bodyReader)
	return Response{tester, tester.PerformRequest(request)}
}

func (tester *Tester) GetSquadList() Response {
	return tester.DoRequest("GET", "/squad", nil)
}

func (tester *Tester) GetSquad(squadId string) Response {
	return tester.DoRequest("GET", "/squad/"+squadId, nil)
}

func (tester *Tester) PostSquad() Response {
	return tester.DoRequest("POST", "/squad", "")
}

func (tester *Tester) PostSquadMember(squadId string, member api.SquadMember) Response {
	return tester.DoRequest("POST", "/squad/"+squadId, member)
}

func (tester *Tester) PerformPostSquad() string {
	var newSquadId string
	tester.PostSquad().
		CheckStatus(http.StatusAccepted).
		LoadJson(&newSquadId)
	return newSquadId
}

func (tester *Tester) PerformGetSquad(squadId string) api.Squad {
	squad := api.Squad{}
	tester.GetSquad(squadId).
		CheckStatus(http.StatusOK).
		LoadJson(&squad)
	return squad
}

func (tester *Tester) PerformGetSquadList() []string {
	var loadedJson []string
	tester.GetSquadList().
		CheckStatus(http.StatusOK).
		LoadJson(&loadedJson)
	return loadedJson
}

func (tester *Tester) PerformPostSquadMember(squadId string, squadMember api.SquadMember) string {
	var newSquadMemberId string
	tester.PostSquadMember(squadId, squadMember).
		CheckStatus(http.StatusAccepted).
		LoadJson(&newSquadMemberId)
	return newSquadMemberId
}

type Response struct {
	Tester   *Tester
	Recorder *httptest.ResponseRecorder
}

func (response Response) CheckStatus(status int) Response {
	assert.Equal(response.Tester.t, status, response.Recorder.Code)
	return response
}

func (response Response) LoadJson(loadLocation interface{}) Response {
	err := json.Unmarshal(response.Recorder.Body.Bytes(), loadLocation)
	if err != nil {
		response.Tester.t.Fatal(err)
	}
	return response
}

func newRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
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
