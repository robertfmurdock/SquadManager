package service

import (
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SquadRepositoryFactory struct {
	Config        Configuration
	parentSession *mgo.Session
}

func (factory SquadRepositoryFactory) Close() {
	if factory.parentSession != nil {
		factory.parentSession.Close()
	}
}

func (factory *SquadRepositoryFactory) Repository() (*SquadRepository, error) {
	if factory.parentSession == nil {
		if err := factory.initParentSession(); err != nil {
			return nil, err
		}
	}

	repository := SquadRepository{Config: factory.Config, session: factory.parentSession.Copy()}
	return &repository, nil
}

func (factory *SquadRepositoryFactory) initParentSession() error {
	session, err := mgo.DialWithTimeout(factory.Config.Host, factory.Config.DbTimeout)
	factory.parentSession = session
	return err
}

type SquadRepository struct {
	Config  Configuration
	session *mgo.Session
}

func (repository SquadRepository) Close() {
	repository.session.Close()
}

func (repository SquadRepository) Database() *mgo.Database {
	return repository.session.DB(repository.Config.DatabaseName)
}

func (repository SquadRepository) SquadCollection() *mgo.Collection {
	return repository.Database().C("squad")
}

func (repository SquadRepository) SquadMemberCollection() *mgo.Collection {
	return repository.Database().C("squadMember")
}

func (repository SquadRepository) addSquad() (bson.ObjectId, error) {
	id := bson.NewObjectId()

	collection := repository.SquadCollection()
	return id, collection.Insert(SquadDocument{id})
}

func (repository *SquadRepository) getSquad(idString string) (*api.Squad, error) {

	if !bson.IsObjectIdHex(idString) {
		return nil, nil
	}

	squadId := bson.ObjectIdHex(idString)

	squadCollection := repository.SquadCollection()
	var squadDocuments []SquadDocument

	if err := squadCollection.FindId(squadId).All(&squadDocuments); err != nil {
		return nil, err
	}

	if len(squadDocuments) == 0 {
		return nil, nil
	}

	squadMemberCollection := repository.SquadMemberCollection()

	squadMemberDocuments := []SquadMemberDocument{}

	err := squadMemberCollection.Find(bson.M{"squadId": squadId}).All(&squadMemberDocuments)
	if err != nil {
		return nil, err
	}
	return &api.Squad{
		ID:      idString,
		Members: toApiSquadMemberList(squadMemberDocuments),
	}, nil
}

func toApiSquadMemberList(documents []SquadMemberDocument) []api.SquadMember {
	idList := make([]api.SquadMember, len(documents))
	for index, document := range documents {
		idList[index] = toApiSquadMember(document)
	}
	return idList
}

func toApiSquadMember(document SquadMemberDocument) api.SquadMember {
	return api.SquadMember{
		ID:    document.ID.Hex(),
		Range: document.Range,
		Email: document.Email,
	}
}

func (repository SquadRepository) postSquadMember(squadMember api.SquadMember, squadId string) error {
	collection := repository.SquadMemberCollection()
	return collection.Insert(toSquadMemberDocument(squadMember, bson.ObjectIdHex(squadId)))
}

func toSquadMemberDocument(squadMember api.SquadMember, squadId bson.ObjectId) SquadMemberDocument {
	return SquadMemberDocument{
		ID:      bson.ObjectIdHex(squadMember.ID),
		Email:   squadMember.Email,
		Range:   squadMember.Range,
		SquadID: squadId,
	}
}

func (repository SquadRepository) findSquadDocuments(query interface{}) ([]SquadDocument, error) {
	collection := repository.SquadCollection()
	var squadDocuments []SquadDocument
	err := collection.Find(query).All(&squadDocuments)
	return squadDocuments, err
}

func (repository SquadRepository) listSquads() ([]string, error) {

	squadDocuments, err := repository.findSquadDocuments(bson.M{})

	if err != nil {
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

type SquadMemberDocument struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	SquadID bson.ObjectId `bson:"squadId"`
	Range   api.Range     `bson:"range"`
	Email   string        `bson:"email"`
}
