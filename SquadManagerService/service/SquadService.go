package service

import (
	"github.com/julienschmidt/httprouter"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"encoding/json"
	"net/http"
)

type SquadService struct {
	RepositoryFactory *SquadRepositoryFactory
}

func (self SquadService) Close() {
	self.RepositoryFactory.Close()
}

func newSquadService(config Configuration) (*SquadService, error) {
	repositoryFactory := SquadRepositoryFactory{config, nil}

	squadService := SquadService{&repositoryFactory}

	return &squadService, nil
}

func (self SquadService) listSquads(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	repository, err := self.RepositoryFactory.Repository()
	if err != nil {
		writer.WriteHeader(500)
		return
	}
	defer repository.Close()

	squads, err := repository.listSquads()
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	json.NewEncoder(writer).Encode(squads)
}

func (self SquadService) createSquad(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	repository, err := self.RepositoryFactory.Repository()
	defer repository.Close()
	squadId, err := repository.addSquad()
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadId)
}

func (self SquadService) getSquad(writer http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	repository, err := self.RepositoryFactory.Repository()
	defer repository.Close()
	id := params.ByName("id")
	squad, err := repository.getSquad(id)
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
	repository, _ := self.RepositoryFactory.Repository()

	if err := repository.postSquadMember(squadMember, squadId); err != nil {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadMember.ID)
}
