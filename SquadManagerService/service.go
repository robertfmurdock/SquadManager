package service

import (
	"net/http"
	"encoding/json"
)

func MainHandler(writer http.ResponseWriter, request *http.Request) {

	squads := []string{}

	json.NewEncoder(writer).Encode(squads)
}

