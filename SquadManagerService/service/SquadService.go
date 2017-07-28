package service

import (
	"github.com/julienschmidt/httprouter"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"encoding/json"
	"net/http"
)

type SquadService struct {
	Repository *SquadRepository
}

func newSquadService(config Configuration) (*SquadService, error) {
	repository, err := newSquadRepository(config)
	if err != nil {
		return nil, err
	}

	squadService := SquadService{repository}

	return &squadService, nil
}

func (self SquadService) listSquads(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squads, err := self.Repository.listSquads()
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	json.NewEncoder(writer).Encode(squads)
}

func (self SquadService) createSquad(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squadId, err := self.Repository.addSquad()
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadId)
}

func (self SquadService) getSquad(writer http.ResponseWriter, _ *http.Request, params httprouter.Params) {

	id := params.ByName("id")
	squad, err := self.Repository.getSquad(id)
	if err != nil {
		writer.WriteHeader(500)
	}

	json.NewEncoder(writer).Encode(squad)
}

func (self SquadService) postSquadMember(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	squadId := params.ByName("id")
	var squadMember api.SquadMember

	if json.NewDecoder(request.Body).Decode(&squadMember) != nil {
		writer.WriteHeader(500)
		return
	}

	err := self.Repository.postSquadMember(squadMember, squadId)
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadMember.ID)
}
