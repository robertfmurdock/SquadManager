package service_test

import (
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/testutility"
	"gopkg.in/mgo.v2/bson"
)

var (
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "localhost",
	}
	mainHandler = service.MakeMainHandler(config)
)

func TestNoResponseOnMainUrl(t *testing.T) {
	request := testutility.NewRequest(t, "GET", "/", nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 404)
}

func TestPOSTSquadWillIncludeNewSquadInSubsequentGET(t *testing.T) {
	wrapper := testutility.Wrap(t, mainHandler)

	newSquadId := wrapper.PerformPostSquadAndGetId()
	squadList := wrapper.PerformGetSquadList(t)

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

	member := api.SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Email: "fakeemail@fake.com",
	}
	memberId := wrapper.PerformPostSquadMember(newSquadId, member)

	assert.Equal(t, member.ID, memberId)

	squad := wrapper.PerformGetSquad(newSquadId)

	assert.Contains(t, squad.Members, member)
}
