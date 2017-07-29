package service

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"time"
)

type Configuration struct {
	DatabaseName string
	Host         string
	DbTimeout    time.Duration
}

func MakeMainHandler(config Configuration) http.Handler {

	repository, err := newSquadService(config)
	if err != nil {
		panic(err)
	}
	defer repository.Close()

	router := httprouter.New()

	router.GET("/squad", repository.listSquads)
	router.POST("/squad", repository.createSquad)
	router.GET("/squad/:id", repository.getSquad)
	router.POST("/squad/:id", repository.postSquadMember)

	return http.HandlerFunc(router.ServeHTTP)
}
