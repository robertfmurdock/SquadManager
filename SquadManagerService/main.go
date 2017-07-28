package main

import (
	"github.com/robertfmurdock/SquadManager/SquadManagerService/service"
	"github.com/urfave/negroni"
)

func main() {
	handler := service.MakeMainHandler(service.Configuration{})

	classic := negroni.Classic()
	classic.UseHandler(handler)

	classic.Run(":8080")
}