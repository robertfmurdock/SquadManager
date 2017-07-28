package service

import (
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
)

type SquadService struct {
	Repository *SquadRepository
}

func NewSquadService(config ServiceConfiguration) (*SquadService, error) {
	repository, err := NewSquadRepository(config)
	if (err != nil) {
		return nil, err
	}

	service := SquadService{repository}

	return &service, nil
}

func (self SquadService) listSquads(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squads, err := self.Repository.listSquads()
	if (err != nil) {
		writer.WriteHeader(500)
		return
	}

	json.NewEncoder(writer).Encode(squads)
}

func (self SquadService) createSquad(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squadId, err := self.Repository.addSquad()
	if (err != nil) {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadId)
}

