package service

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ResponseEntity struct {
	value interface{}
	code  int
}

type ThinHandler func(_ *http.Request, _ httprouter.Params) (ResponseEntity, error)

func (handler ThinHandler) With(service *Context) httprouter.Handle {

	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		entity, err := handler(request, params)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.WriteHeader(entity.code)
		json.NewEncoder(writer).Encode(entity.value)
	}
}

type Handler func(_ *http.Request, _ httprouter.Params, _ *SquadRepository) (ResponseEntity, error)

func (handler Handler) With(service *Context) httprouter.Handle {

	return ThinHandler(func(
		request *http.Request,
		params httprouter.Params,
	) (ResponseEntity, error) {

		repository, err := service.RepositoryFactory.Repository()
		if err != nil {
			return ResponseEntity{}, err
		}
		defer repository.Close()
		return handler(request, params, repository)

	}).With(service)
}

type NoInputHandler func(_ *SquadRepository) (ResponseEntity, error)

func (handler NoInputHandler) With(service *Context) httprouter.Handle {
	return Handler(func(
		request *http.Request,
		params httprouter.Params,
		repository *SquadRepository,
	) (ResponseEntity, error) {
		return handler(repository)
	}).With(service)
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
