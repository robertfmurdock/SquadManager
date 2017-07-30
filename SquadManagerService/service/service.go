package service

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Configuration struct {
	DatabaseName string
	Host         string
	DbTimeout    time.Duration
}

func MakeMainHandler(config Configuration) http.Handler {

	context, err := newContext(config)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	router := httprouter.New()

	router.GET("/squad", context.with(NoInputHandler(listSquads)))
	router.POST("/squad", context.with(NoInputHandler(createSquad)))
	router.GET("/squad/:id", context.with(SquadHandler(getSquad)))
	router.POST("/squad/:id", context.with(SquadHandler(postSquadMember)))

	return http.HandlerFunc(router.ServeHTTP)
}

