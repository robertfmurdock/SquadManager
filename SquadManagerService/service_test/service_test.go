package service_test

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"encoding/json"
	"io"
	"bytes"
	"gopkg.in/mgo.v2/bson"
)

var (
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "localhost",
	}
	mainHandler = service.MakeMainHandler(config)
)

func getPostBody(body interface{}) ([]byte, error) {
	if value, ok := body.(string); ok {
		return []byte(value), nil
	} else {
		return json.Marshal(body)
	}
}

func makePostRequest(t *testing.T, url string, body interface{}) *http.Request {
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

func TestNoResponseOnMainUrl(t *testing.T) {
	request := newRequest(t, "GET", "/", nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 404)
}

func newRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func performGetSquadList(t *testing.T) []string {
	request := newRequest(t, "GET", "/squad", nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 200)

	var loadedJson []string
	json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	return loadedJson
}

func performPostSquad(t *testing.T) *httptest.ResponseRecorder {
	request := makePostRequest(t, "/squad", "")

	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)
	return recorder
}

func performPostSquadMember(t *testing.T, squadId string, squadMember api.SquadMember) string {
	request := makePostRequest(t, "/squad/"+squadId, squadMember)

	response := httptest.NewRecorder()

	mainHandler.ServeHTTP(response, request)

	assert.Equal(t, http.StatusAccepted, response.Code)
	var newSquadMemberId string
	json.Unmarshal(response.Body.Bytes(), &newSquadMemberId)

	return newSquadMemberId
}

func performGetSquad(t *testing.T, squadId string) api.Squad {
	request := newRequest(t, "GET", "/squad/"+squadId, nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 200)

	var loadedJson api.Squad
	err := json.Unmarshal(recorder.Body.Bytes(), &loadedJson)

	if err != nil {
		t.Fatal(err)
	}

	return loadedJson
}

func performPostSquadAndGetId(t *testing.T) string {
	response := performPostSquad(t)
	assert.Equal(t, http.StatusAccepted, response.Code)
	var newSquadId string
	json.Unmarshal(response.Body.Bytes(), &newSquadId)
	return newSquadId
}

func TestPOSTSquadWillIncludeNewSquadInSubsequentGET(t *testing.T) {

	newSquadId := performPostSquadAndGetId(t)

	squadList := performGetSquadList(t)

	assert.Contains(t, squadList, newSquadId)
}

func TestGETSquadWithNewSquadWillHaveNoMembers(t *testing.T) {
	newSquadId := performPostSquadAndGetId(t)

	squad := performGetSquad(t, newSquadId)

	expectedSquad := api.Squad{
		ID:      newSquadId,
		Members: []api.SquadMember{},
	}

	assert.Equal(t, expectedSquad, squad)
}

func TestPOSTSquadMemberWillShowSquadMemberInSubsequentGET(t *testing.T) {
	newSquadId := performPostSquadAndGetId(t)

	member := api.SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Email: "fakeemail@fake.com",
	}
	memberId := performPostSquadMember(t, newSquadId, member)

	assert.Equal(t, member.ID, memberId)

	squad := performGetSquad(t, newSquadId)

	assert.Contains(t, squad.Members, member)
}
