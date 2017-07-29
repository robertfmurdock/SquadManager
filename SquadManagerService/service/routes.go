package service

import (
	"encoding/json"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"net/http"
)

func listSquads(repository *SquadRepository) (ResponseEntity, error) {
	squads, err := repository.listSquads()
	return ResponseEntity{squads, http.StatusOK}, err
}

func createSquad(repository *SquadRepository) (ResponseEntity, error) {
	squad, err := repository.addSquad()
	return ResponseEntity{squad, http.StatusAccepted}, err
}

func getSquad(_ *http.Request, repository *SquadRepository, squadId string) (ResponseEntity, error) {
	squad, err := repository.getSquad(squadId)
	return ResponseEntity{squad, http.StatusOK}, err
}

func postSquadMember(request *http.Request, repository *SquadRepository, squadId string) (ResponseEntity, error) {
	var squadMember api.SquadMember

	if err := json.NewDecoder(request.Body).Decode(&squadMember); err != nil {
		return ResponseEntity{}, err
	}

	err := repository.postSquadMember(squadMember, squadId)
	return ResponseEntity{squadMember.ID, http.StatusAccepted}, err
}
