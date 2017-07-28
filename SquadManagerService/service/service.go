package service

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type Configuration struct {
	DatabaseName string
	Host         string
}

func MakeMainHandler(config Configuration) http.Handler {

	repository, err := newSquadService(config)

	if err != nil {
		panic(err)
	}

	router := httprouter.New()

	router.GET("/squad", repository.listSquads)
	router.POST("/squad", repository.createSquad)
	router.GET("/squad/:id", repository.getSquad)
	router.POST("/squad/:id", repository.postSquadMember)


	return http.HandlerFunc(router.ServeHTTP)
}

