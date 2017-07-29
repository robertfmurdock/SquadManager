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

	context, err := newContext(config)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	router := httprouter.New()

	router.GET("/squad", context.with(Handler(listSquads)))
	router.POST("/squad", context.with(Handler(createSquad)))
	router.GET("/squad/:id", context.with(SquadHandler(getSquad)))
	router.POST("/squad/:id", context.with(SquadHandler(postSquadMember)))

	return http.HandlerFunc(router.ServeHTTP)
}

