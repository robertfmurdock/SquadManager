package service

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
)

func listSquads(request *http.Request, _ httprouter.Params, repository *SquadRepository) (ResponseEntity, error) {
	squadParameters, err := parseSquadParameters(request)
	squads, err := repository.listSquads(squadParameters.begin, squadParameters.end)
	return ResponseEntity{squads, http.StatusOK}, err
}

func overwriteSquadList(request *http.Request, _ httprouter.Params, repository *SquadRepository) (ResponseEntity, error) {
	squadList := []api.Squad{}
	if err := json.NewDecoder(request.Body).Decode(&squadList); err != nil {
		return ResponseEntity{err, http.StatusBadRequest}, nil
	}

	squads, err := repository.overwriteSquadList(squadList)
	return ResponseEntity{squads, http.StatusOK}, err
}

func createSquad(repository *SquadRepository) (ResponseEntity, error) {
	squad, err := repository.addSquad()
	return ResponseEntity{squad, http.StatusAccepted}, err
}

type SquadParameters struct {
	begin *time.Time
	end   *time.Time
}

func getSquad(request *http.Request, repository *SquadRepository, squadId string) (ResponseEntity, error) {
	parameters, err := parseSquadParameters(request)
	if err != nil {
		return ResponseEntity{err, http.StatusBadRequest}, nil
	}

	squad, err := repository.getSquad(squadId, parameters.begin, parameters.end)
	if err != nil {
		return ResponseEntity{}, err
	}

	if squad == nil {
		return ResponseEntity{code: http.StatusNotFound}, nil
	}

	return ResponseEntity{squad, http.StatusOK}, nil
}

func parseSquadParameters(request *http.Request) (SquadParameters, error) {
	values := request.URL.Query()
	beginDate, err := api.ParseDate(values.Get("begin"))
	if err != nil {
		return SquadParameters{}, err
	}
	endDate, err := api.ParseDate(values.Get("end"))

	return SquadParameters{beginDate, endDate}, err
}

func postSquadMember(request *http.Request, repository *SquadRepository, squadId string) (ResponseEntity, error) {
	var squadMember api.SquadMember
	if err := json.NewDecoder(request.Body).Decode(&squadMember); err != nil {
		return ResponseEntity{err, http.StatusBadRequest}, nil
	}

	if squad, err := repository.getSquad(squadId, nil, nil); err != nil || squad == nil {
		return ResponseEntity{code: http.StatusNotFound}, err
	}

	err := repository.postSquadMember(squadMember, squadId)
	return ResponseEntity{squadMember.ID, http.StatusAccepted}, err
}
