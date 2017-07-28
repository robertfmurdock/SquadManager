package service

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SquadRepository struct {
	Config  ServiceConfiguration
	session *mgo.Session
}

func NewSquadRepository(config ServiceConfiguration) (*SquadRepository, error) {
	session, err := mgo.Dial("mongo,localhost")
	if err != nil {
		return nil, err
	}

	return &SquadRepository{config, session}, nil
}

func (self SquadRepository) Database() *mgo.Database {
	return self.session.DB(self.Config.DatabaseName)
}

func (self SquadRepository) SquadCollection() *mgo.Collection {
	return self.Database().C("squad")
}

func (self SquadRepository) addSquad() (bson.ObjectId, error) {
	collection := self.SquadCollection()
	id := bson.NewObjectId()
	return id, collection.Insert(SquadDocument{id})
}

func (self SquadRepository) findSquadDocuments(query interface{}) ([]SquadDocument, error) {
	collection := self.SquadCollection()
	var squadDocuments []SquadDocument
	err := collection.Find(query).All(&squadDocuments)
	return squadDocuments, err
}

func (self SquadRepository) listSquads() ([]string, error) {
	squadDocuments, err := self.findSquadDocuments(bson.M{})

	if (err != nil) {
		return nil, err
	}

	return toSquadIds(squadDocuments), nil
}

func toSquadIds(documents []SquadDocument) []string {
	idList := make([]string, len(documents))
	for index, document := range documents {
		idList[index] = document.ID.Hex()
	}
	return idList
}

type SquadDocument struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
}