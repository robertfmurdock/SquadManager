package service_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"net/http"

	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	tu "github.com/robertfmurdock/SquadManager/SquadManagerService/testutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

var (
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "localhost",
		DbTimeout:    time.Second / 100,
	}
	mainHandler = service.MakeMainHandler(config)
)

func TestNoResponseOnMainUrl(t *testing.T) {
	request := tu.NewRequest(t, "GET", "/", nil)
	recorder := httptest.NewRecorder()

	mainHandler.ServeHTTP(recorder, request)

	assert.Equal(t, recorder.Code, 404)
}

func TestWillErrorWhenDatasourceNotAvailable(t *testing.T) {
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "missing",
		DbTimeout:    time.Millisecond / 100,
	}
	handler := service.MakeMainHandler(config)

	wrapper := tu.Wrap(t, handler)

	tu.GetSquadList().
		WithStatus(http.StatusInternalServerError).
		Perform(wrapper)
}

func TestPOSTSquadWillIncludeNewSquadInSubsequentGET(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)

	newSquadId := wrapper.PerformPostSquad()
	squadList := wrapper.PerformGetSquadList()

	assert.Contains(t, squadList, newSquadId)
}

func TestGETSquadWithNewSquadWillHaveNoMembers(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)
	newSquadId := wrapper.PerformPostSquad()

	squad := wrapper.PerformGetSquad(newSquadId)

	expectedSquad := api.Squad{
		ID:      newSquadId,
		Members: []api.SquadMember{},
	}

	assert.Equal(t, expectedSquad, squad)
}

func TestGETSquadWithUnknownSquadIdWillReturn404(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)

	squadId := bson.NewObjectId().Hex()

	tu.GetSquad(squadId).
		WithStatus(http.StatusNotFound).
		Perform(wrapper)
}

func TestGETSquadWithInvalidSquadIdWillReturn404(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)

	squadId := "This is not a valid object id"

	tu.GetSquad(squadId).
		WithStatus(http.StatusNotFound).
		Perform(wrapper)
}

func TestPOSTSquadMemberWillShowSquadMemberInSubsequentGET(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)
	newSquadId := wrapper.PerformPostSquad()

	now := time.Now().Truncate(24 * time.Hour)
	later := now.AddDate(1, 0, 0)

	member := api.SquadMember{
		ID: bson.NewObjectId().Hex(),
		Range: api.Range{
			Begin: now,
			End:   later,
		},
		Email: "fakeemail@fake.com",
	}
	memberId := wrapper.PerformPostSquadMember(newSquadId, member)

	assert.Equal(t, member.ID, memberId)

	squad := wrapper.PerformGetSquad(newSquadId)

	assert.Contains(t, squad.Members, member)
}

func TestPOSTSquadMemberWill404WhenSquadDoesNotExist(t *testing.T) {
	wrapper := tu.Wrap(t, mainHandler)

	now := time.Now().Truncate(24 * time.Hour)
	later := now.AddDate(1, 0, 0)

	member := api.SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Range: api.Range{Begin: now, End: later},
		Email: "fakeemail@fake.com",
	}
	squadId := bson.NewObjectId().Hex()

	tu.PostSquadMember(squadId, member).
		WithStatus(http.StatusNotFound).
		Perform(wrapper)
}
