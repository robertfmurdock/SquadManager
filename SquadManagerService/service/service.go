package service

import (
	"net/http"
	"encoding/json"
)

func MakeMainHandler() http.Handler {
	return http.HandlerFunc(mainHandler)
}

func mainHandler(writer http.ResponseWriter, request *http.Request) {

	mux := http.NewServeMux()

	mux.HandleFunc("/squad", SquadListHandler)

	mux.ServeHTTP(writer, request)
}

func SquadListHandler(writer http.ResponseWriter, _ *http.Request) {

	squads := []string{}
	json.NewEncoder(writer).Encode(squads)
}

