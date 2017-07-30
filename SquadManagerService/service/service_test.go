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

	squad := tester.PerformGetSquad(newSquadId, nil, nil)

	expectedSquad := api.Squad{
		ID:      newSquadId,
		Members: []api.SquadMember{},
	}

	assert.Equal(t, expectedSquad, squad)
}

func TestGETSquadWithUnknownSquadIdWillReturn404(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	squadId := bson.NewObjectId().Hex()

	tester.GetSquad(squadId, nil, nil).
		CheckStatus(http.StatusNotFound)
}

func TestGETSquadWithInvalidSquadIdWillReturn404(t *testing.T) {
	tester := testutil.New(t, mainHandler)

	squadId := "This is not a valid object id"

	tester.GetSquad(squadId, nil, nil).
		CheckStatus(http.StatusNotFound)
}

func TestPOSTSquadMembersWillShowSquadMembersInSubsequentGET(t *testing.T) {
	tester := testutil.New(t, mainHandler)
	newSquadId := tester.PerformPostSquad()

	members := []api.SquadMember{
		api.NewSquadMember("dale@fake.com",
			api.Range{
				Begin: *api.Date(2017, 07, 30),
				End:   *api.Date(2017, 11, 10),
			}),
		api.NewSquadMember("chip@fake.com",
			api.Range{
				Begin: *api.Date(2017, 05, 1),
				End:   *api.Date(2017, 9, 15),
			}),
		api.NewSquadMember("daisy@fake.com",
			api.Range{
				Begin: *api.Date(2017, 11, 15),
				End:   *api.Date(2018, 2, 7),
			}),
	}

	for _, member := range members {
		tester.PerformPostSquadMember(newSquadId, member)
	}

	squad := tester.PerformGetSquad(newSquadId, nil, nil)
	assert.Equal(t, squad.Members, members)
}

func TestPOSTSquadMemberMultipleTimesWillUpdate(t *testing.T) {
	tester := testutil.New(t, mainHandler)
	squadId := tester.PerformPostSquad()

	member := api.NewSquadMember("dale@fake.com",
		api.Range{
			Begin: *api.Date(2017, 07, 30),
			End:   *api.Date(2017, 11, 10),
		})

	tester.PerformPostSquadMember(squadId, member)

	member.Range.End = *api.Date(2017, 8, 10)

	tester.PerformPostSquadMember(squadId, member)

	squad := tester.PerformGetSquad(squadId, nil, nil)
	assert.Equal(t, []api.SquadMember{member}, squad.Members)
}

func TestSquadMembersCanBeFilteredInGET(t *testing.T) {
	tester := testutil.New(t, mainHandler)
	newSquadId := tester.PerformPostSquad()

	members := []api.SquadMember{
		api.NewSquadMember("dale@fake.com",
			api.Range{
				Begin: *api.Date(2017, 7, 30),
				End:   *api.Date(2017, 8, 10),
			}),
		api.NewSquadMember("chip@fake.com",
			api.Range{
				Begin: *api.Date(2017, 8, 11),
				End:   *api.Date(2017, 9, 15),
			}),
		api.NewSquadMember("daisy@fake.com",
			api.Range{
				Begin: *api.Date(2017, 9, 20),
				End:   *api.Date(2018, 2, 7),
			}),
	}

	for _, member := range members {
		tester.PerformPostSquadMember(newSquadId, member)
	}

	squad := tester.PerformGetSquad(newSquadId, api.Date(2017, 8, 11), api.Date(2017, 9, 15))
	assert.Equal(t, members[1:2], squad.Members)
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
