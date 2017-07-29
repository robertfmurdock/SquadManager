package service

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
)

type SquadRepositoryFactory struct {
	Config        Configuration
	parentSession *mgo.Session
}

func (self SquadRepositoryFactory) Close() {
	if self.parentSession != nil {
		self.parentSession.Close()
	}
}

func (self *SquadRepositoryFactory) Repository() (*SquadRepository, error) {
	if self.parentSession == nil {
		if err := self.initParentSession(); err != nil {
			return nil, err
		}
	}

	repository := SquadRepository{Config: self.Config, session: self.parentSession.Copy()}
	return &repository, nil
}

func (self *SquadRepositoryFactory) initParentSession() error {
	session, err := mgo.DialWithTimeout(self.Config.Host, self.Config.DbTimeout)
	self.parentSession = session
	return err
}

type SquadRepository struct {
	Config  Configuration
	session *mgo.Session
}

func (self SquadRepository) Close() {
	self.session.Close()
}

func (self SquadRepository) Database() *mgo.Database {
	return self.session.DB(self.Config.DatabaseName)
}

func (self SquadRepository) SquadCollection() *mgo.Collection {
	return self.Database().C("squad")
}

func (self SquadRepository) SquadMemberCollection() *mgo.Collection {
	return self.Database().C("squadMember")
}

func (self SquadRepository) addSquad() (bson.ObjectId, error) {
	id := bson.NewObjectId()

	collection := self.SquadCollection()
	return id, collection.Insert(SquadDocument{id})
}

func (self *SquadRepository) getSquad(id string) (*api.Squad, error) {
	collection := self.SquadMemberCollection()

	squadDocument := []SquadMemberDocument{}
	err := collection.Find(bson.M{"squadId": bson.ObjectIdHex(id)}).All(&squadDocument)
	if err != nil {
		return nil, err
	}
	return &api.Squad{
		ID:      id,
		Members: toApiSquadMemberList(squadDocument),
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

func (self SquadRepository) postSquadMember(squadMember api.SquadMember, squadId string) (error) {
	collection := self.SquadMemberCollection()
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

func (self SquadRepository) findSquadDocuments(query interface{}) ([]SquadDocument, error) {
	collection := self.SquadCollection()
	var squadDocuments []SquadDocument
	err := collection.Find(query).All(&squadDocuments)
	return squadDocuments, err
}

func (self SquadRepository) listSquads() ([]string, error) {

	squadDocuments, err := self.findSquadDocuments(bson.M{})

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
	Range   api.Range `bson:"range"`
	Email   string `bson:"email"`
}
