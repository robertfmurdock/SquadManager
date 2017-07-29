package service_test

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/testutility"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "localhost",
		DbTimeout: time.Second / 100,
	}
	mainHandler = service.MakeMainHandler(config)
)

func TestNoResponseOnMainUrl(t *testing.T) {
	request := testutility.NewRequest(t, "GET", "/", nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 404)
}

func TestWillErrorWhenDatasourceNotAvailable( t *testing.T) {
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "missing",
		DbTimeout: time.Millisecond / 100,
	}
	handler := service.MakeMainHandler(config)

	wrapper := testutility.Wrap(t, handler)

	request := testutility.NewRequest(t, "GET", "/squad", nil)

	response := wrapper.PerformRequest(request)

	assert.Equal(t, response.Code, 500)

}

func TestPOSTSquadWillIncludeNewSquadInSubsequentGET(t *testing.T) {
	wrapper := testutility.Wrap(t, mainHandler)

	newSquadId := wrapper.PerformPostSquadAndGetId()
	squadList := wrapper.PerformGetSquadList()

	assert.Contains(t, squadList, newSquadId)
}

func TestGETSquadWithNewSquadWillHaveNoMembers(t *testing.T) {
	wrapper := testutility.Wrap(t, mainHandler)
	newSquadId := wrapper.PerformPostSquadAndGetId()

	squad := wrapper.PerformGetSquad(newSquadId)

	expectedSquad := api.Squad{
		ID:      newSquadId,
		Members: []api.SquadMember{},
	}

	assert.Equal(t, expectedSquad, squad)
}

func TestPOSTSquadMemberWillShowSquadMemberInSubsequentGET(t *testing.T) {
	wrapper := testutility.Wrap(t, mainHandler)
	newSquadId := wrapper.PerformPostSquadAndGetId()

	now := time.Now().Truncate(24 * time.Hour)
	later := now.AddDate(1, 0, 0)

	member := api.SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Range: api.Range{Begin: now, End: later},
		Email: "fakeemail@fake.com",
	}
	memberId := wrapper.PerformPostSquadMember(newSquadId, member)

	assert.Equal(t, member.ID, memberId)

	squad := wrapper.PerformGetSquad(newSquadId)

	assert.Contains(t, squad.Members, member)
}
