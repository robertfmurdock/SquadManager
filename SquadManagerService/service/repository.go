package service

import (
	"time"

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

func (repository *SquadRepository) getSquad(idString string, begin *time.Time, end *time.Time) (*api.Squad, error) {
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

	return repository.loadSquad(squadId, begin, end)
}

func (repository *SquadRepository) loadSquad(squadId bson.ObjectId, begin *time.Time, end *time.Time) (*api.Squad, error) {
	squadMemberDocuments := []SquadMemberDocument{}
	if err := repository.loadSquadMemberDocuments(bson.M{"squadId": squadId}, &squadMemberDocuments); err != nil {
		return nil, err
	}

	return buildSquad(squadId, squadMemberDocuments, begin, end), nil
}

func (repository *SquadRepository) loadSquadMemberDocuments(query interface{}, squadMemberDocuments *[]SquadMemberDocument) error {
	return repository.SquadMemberCollection().Find(query).All(squadMemberDocuments)
}

func buildSquad(squadId bson.ObjectId, squadMemberDocuments []SquadMemberDocument, begin *time.Time, end *time.Time) *api.Squad {
	squad := &api.Squad{
		ID:      squadId.Hex(),
		Members: api.FilterMembers(toApiSquadMemberList(squadMemberDocuments), begin, end),
	}
	return squad
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
	squadMemberDocument := toSquadMemberDocument(squadMember, bson.ObjectIdHex(squadId))
	_, err := collection.Upsert(bson.M{"_id": squadMemberDocument.ID}, squadMemberDocument)
	return err
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

func (repository SquadRepository) listSquads() ([]api.Squad, error) {
	squadDocuments, err := repository.findSquadDocuments(bson.M{})

	if err != nil {
		return nil, err
	}

	var allSquadMemberDocuments []SquadMemberDocument
	repository.loadSquadMemberDocuments(bson.M{}, &allSquadMemberDocuments)

	squadList := make([]api.Squad, len(squadDocuments))
	for index, document := range squadDocuments {
		relatedSquadMemberDocuments := filterBySquadId(document.ID, allSquadMemberDocuments)
		squad := buildSquad(document.ID, relatedSquadMemberDocuments, nil, nil)

		squadList[index] = *squad
	}

	return squadList, nil
}
func filterBySquadId(squadId bson.ObjectId, documents []SquadMemberDocument) []SquadMemberDocument {
	var results []SquadMemberDocument
	for _, document := range documents {
		if document.SquadID == squadId {
			results = append(results, document)
		}
	}
	return results
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
