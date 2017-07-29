package service

import (
	"github.com/julienschmidt/httprouter"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"encoding/json"
	"net/http"
)

type Context struct {
	RepositoryFactory *SquadRepositoryFactory
}

func newContext(config Configuration) (*Context, error) {
	repositoryFactory := SquadRepositoryFactory{config, nil}

	squadService := Context{&repositoryFactory}

	return &squadService, nil
}

type ContextualHandler interface {
	With(service *Context) httprouter.Handle
}

func (context *Context) with(contextual ContextualHandler) httprouter.Handle {
	return contextual.With(context)
}

func (context Context) Close() {
	context.RepositoryFactory.Close()
}

type ResponseEntity struct {
	value interface{}
	code  int
}

type Handler func(_ *http.Request, _ httprouter.Params, _ *SquadRepository) (ResponseEntity, error)

func (handler Handler) With(service *Context) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		repository, err := service.RepositoryFactory.Repository()
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		defer repository.Close()

		responseEntity, err2 := handler(request, params, repository)
		if err2 != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(responseEntity.code)
		json.NewEncoder(writer).Encode(responseEntity.value)
	}
}

type SquadHandler func(_ *http.Request, _ *SquadRepository, _ string) (ResponseEntity, error)

func (handler SquadHandler) With(service *Context) httprouter.Handle {
	return Handler(func(
		request *http.Request,
		params httprouter.Params,
		repository *SquadRepository,
	) (ResponseEntity, error) {
		squadId := params.ByName("id")
		return handler(request, repository, squadId)
	}).With(service)
}

func listSquads(_ *http.Request, _ httprouter.Params, repository *SquadRepository) (ResponseEntity, error) {
	squads, err := repository.listSquads()
	return ResponseEntity{squads, http.StatusOK}, err
}

func createSquad(_ *http.Request, _ httprouter.Params, repository *SquadRepository) (ResponseEntity, error) {
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
