package service

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type ServiceConfiguration struct {
	DatabaseName string
	Host         string
}

func MakeMainHandler(config ServiceConfiguration) http.Handler {

	repository, err := NewSquadService(config)

	if (err != nil) {
		panic(err)
	}

	router := httprouter.New()

	router.GET("/squad", repository.listSquads)
	router.POST("/squad", repository.createSquad)

	return http.HandlerFunc(router.ServeHTTP)
}

