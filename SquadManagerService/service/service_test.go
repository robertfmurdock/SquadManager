package service_test

import (
	"testing"
	"time"

	"net/http"

	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/testutil"
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
	tester := testutil.New(t, mainHandler)

	tester.DoRequest("GET", "/", nil).
		CheckStatus(http.StatusNotFound)
}

func TestWillErrorWhenDatasourceNotAvailable(t *testing.T) {
	config = service.Configuration{
		DatabaseName: "SquadManagerTestDB",
		Host:         "missing",
		DbTimeout:    time.Millisecond / 100,
	}
	handler := service.MakeMainHandler(config)

	tester := testutil.New(t, handler)

	tester.GetSquadList().
		CheckStatus(http.StatusInternalServerError)
}

func TestPOSTSquadWillIncludeNewSquadInSubsequentGET(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	newSquadId := tester.PerformPostSquad()
	squadList := tester.PerformGetSquadList()

	assert.Contains(t, squadList, newSquadId)
}

func TestGETSquadWithNewSquadWillHaveNoMembers(t *testing.T) {
	tester := testutil.New(t, mainHandler)
	newSquadId := tester.PerformPostSquad()

	squad := tester.PerformGetSquad(newSquadId)

	expectedSquad := api.Squad{
		ID:      newSquadId,
		Members: []api.SquadMember{},
	}

	assert.Equal(t, expectedSquad, squad)
}

func TestGETSquadWithUnknownSquadIdWillReturn404(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	squadId := bson.NewObjectId().Hex()

	tester.GetSquad(squadId).
		CheckStatus(http.StatusNotFound)
}

func TestGETSquadWithInvalidSquadIdWillReturn404(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	squadId := "This is not a valid object id"

	tester.GetSquad(squadId).
		CheckStatus(http.StatusNotFound)
}

func TestPOSTSquadMemberWillShowSquadMemberInSubsequentGET(t *testing.T) {
	tester := testutil.New(t, mainHandler)
	newSquadId := tester.PerformPostSquad()

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

	memberId := tester.PerformPostSquadMember(newSquadId, member)
	assert.Equal(t, member.ID, memberId)

	squad := tester.PerformGetSquad(newSquadId)
	assert.Contains(t, squad.Members, member)
}

func TestPOSTSquadMemberWill404WhenSquadDoesNotExist(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	now := time.Now().Truncate(24 * time.Hour)
	later := now.AddDate(1, 0, 0)

	member := api.SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Range: api.Range{Begin: now, End: later},
		Email: "fakeemail@fake.com",
	}
	squadId := bson.NewObjectId().Hex()

	tester.PostSquadMember(squadId, member).
		CheckStatus(http.StatusNotFound)
}
