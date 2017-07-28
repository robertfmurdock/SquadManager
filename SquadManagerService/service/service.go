package service

import (
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ServiceConfiguration struct {
	DatabaseName string
	Host         string
}

type SquadRepository struct {
	Config  ServiceConfiguration
	session *mgo.Session
}

func NewSquadRepository(config ServiceConfiguration) (*SquadRepository, error) {
	session, err := mgo.Dial("mongo,localhost")
	if err != nil {
		return nil, err
	}

	repository := SquadRepository{config, session}

	return &repository, nil
}

func MakeMainHandler(config ServiceConfiguration) http.Handler {

	repository, err := NewSquadRepository(config)

	if (err != nil) {
		panic(err)
	}

	router := httprouter.New()

	router.GET("/squad", repository.squadListHandler)
	router.POST("/squad", repository.squadCreateHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

func (self SquadRepository) squadListHandler(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squads, err := self.listSquads()
	if(err != nil) {
		writer.WriteHeader(500)
		return
	}

	json.NewEncoder(writer).Encode(squads)
}

func (self SquadRepository) squadCreateHandler(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	squadId, err := self.addSquad()
	if (err != nil) {
		writer.WriteHeader(500)
		return
	}

	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(squadId)
}

type SquadDocument struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
}

func (self SquadRepository) addSquad() (bson.ObjectId, error) {
	collection := self.session.DB(self.Config.DatabaseName).C("squad")
	id := bson.NewObjectId()
	return id, collection.Insert(SquadDocument{id})
}

func (self SquadRepository) listSquads() ([]string, error) {
	collection := self.session.DB(self.Config.DatabaseName).C("squad")

	var squadDocuments []SquadDocument
	err := collection.Find(bson.M{}).All(&squadDocuments)
	if(err != nil) {
		return nil, err
	}

	return convertSquadDocumentsToIds(squadDocuments), nil
}

func convertSquadDocumentsToIds(documents []SquadDocument) []string {
	idList := make([]string, len(documents))
	for index, document := range documents {
		idList[index] = document.ID.Hex()
	}
	return idList
}