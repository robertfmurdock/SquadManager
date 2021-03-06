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

func (factory *SquadRepositoryFactory) Close() {
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

func (repository SquadRepository) addSquad() (api.SquadId, error) {
	id := api.SquadId(bson.NewObjectId())
	collection := repository.SquadCollection()
	return id, collection.Insert(SquadDocument{bson.ObjectId(id)})
}

func (repository SquadRepository) overwriteSquadList(squadList []api.Squad) ([]api.Squad, error) {

	if err := clearCollection(repository.SquadCollection()); err != nil {
		return nil, err
	}
	if err := clearCollection(repository.SquadMemberCollection()); err != nil {
		return nil, err
	}

	squadDocumentList, squadMemberDocumentList := toDocuments(squadList)

	if err := insertDocuments(repository.SquadCollection(), squadDocumentList); err != nil {
		return nil, err
	}

	if err := insertDocuments(repository.SquadMemberCollection(), squadMemberDocumentList); err != nil {
		return nil, err
	}

	return squadList, nil
}

func toDocuments(squadList []api.Squad) ([]interface{}, []interface{}) {
	squadDocumentList := make([]interface{}, len(squadList))
	var squadMemberDocumentList []interface{}
	for index, squad := range squadList {
		squadDocumentList[index] = SquadDocument{ID: bson.ObjectId(squad.ID)}

		for _, member := range squad.Members {
			squadMemberDocumentList = append(squadMemberDocumentList, toSquadMemberDocument(member, squad.ID))
		}
	}
	return squadDocumentList, squadMemberDocumentList
}

func clearCollection(collection *mgo.Collection) error {
	count, err := collection.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	return collection.DropCollection()
}

func insertDocuments(collection *mgo.Collection, documentList []interface{}) error {
	if len(documentList) == 0 {
		return nil
	}

	if err := collection.Insert(documentList...); err != nil {
		return err
	}

	return nil
}

func (repository *SquadRepository) getSquad(idString string, begin *time.Time, end *time.Time) (*api.Squad, error) {
	if !bson.IsObjectIdHex(idString) {
		return nil, nil
	}

	squadId := api.SquadId(bson.ObjectIdHex(idString))

	squadCollection := repository.SquadCollection()
	var squadDocuments []SquadDocument

	if err := squadCollection.FindId(bson.ObjectId(squadId)).All(&squadDocuments); err != nil {
		return nil, err
	}

	if len(squadDocuments) == 0 {
		return nil, nil
	}

	return repository.loadSquad(squadId, begin, end)
}

func (repository *SquadRepository) loadSquad(squadId api.SquadId, begin *time.Time, end *time.Time) (*api.Squad, error) {
	query := bson.M{"squadId": bson.ObjectId(squadId)}
	squadMemberDocuments := []SquadMemberDocument{}
	if err := repository.loadSquadMemberDocuments(query, &squadMemberDocuments); err != nil {
		return nil, err
	}
	return buildSquad(squadId, squadMemberDocuments, begin, end), nil
}

func (repository *SquadRepository) loadSquadMemberDocuments(query interface{}, squadMemberDocuments *[]SquadMemberDocument) error {
	return repository.SquadMemberCollection().Find(query).All(squadMemberDocuments)
}

func buildSquad(squadId api.SquadId, squadMemberDocuments []SquadMemberDocument, begin *time.Time, end *time.Time) *api.Squad {
	squad := &api.Squad{
		ID:      squadId,
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
		ID: api.SquadMemberId(document.ID),
		Range: api.Range{
			Begin: document.Range.Begin.UTC(),
			End:   document.Range.End.UTC(),
		},
		Email: document.Email,
	}
}

func (repository SquadRepository) postSquadMember(squadMember api.SquadMember, squadId string) error {
	collection := repository.SquadMemberCollection()
	squadMemberDocument := toSquadMemberDocument(squadMember, api.SquadId(bson.ObjectIdHex(squadId)))
	_, err := collection.Upsert(bson.M{"_id": squadMemberDocument.ID}, squadMemberDocument)
	return err
}

func toSquadMemberDocument(squadMember api.SquadMember, squadId api.SquadId) SquadMemberDocument {
	return SquadMemberDocument{
		ID:      bson.ObjectId(squadMember.ID),
		SquadID: bson.ObjectId(squadId),
		Email:   squadMember.Email,
		Range:   squadMember.Range,
	}
}

func (repository SquadRepository) findSquadDocuments(query interface{}) ([]SquadDocument, error) {
	collection := repository.SquadCollection()
	var squadDocuments []SquadDocument
	err := collection.Find(query).All(&squadDocuments)
	return squadDocuments, err
}

func (repository SquadRepository) listSquads(begin *time.Time, end *time.Time) ([]api.Squad, error) {
	squadDocuments, err := repository.findSquadDocuments(bson.M{})

	if err != nil {
		return nil, err
	}

	var allSquadMemberDocuments []SquadMemberDocument
	repository.loadSquadMemberDocuments(bson.M{}, &allSquadMemberDocuments)

	squadList := []api.Squad{}
	for _, document := range squadDocuments {
		squadId := api.SquadId(document.ID)
		relatedSquadMemberDocuments := filterBySquadId(squadId, allSquadMemberDocuments)
		squad := buildSquad(squadId, relatedSquadMemberDocuments, begin, end)

		noRangeRestrictions := begin == nil && end == nil
		if noRangeRestrictions || len(squad.Members) != 0 {
			squadList = append(squadList, *squad)
		}
	}

	return squadList, nil
}
func filterBySquadId(squadId api.SquadId, documents []SquadMemberDocument) []SquadMemberDocument {
	var results []SquadMemberDocument
	for _, document := range documents {
		if api.SquadId(document.SquadID) == squadId {
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
