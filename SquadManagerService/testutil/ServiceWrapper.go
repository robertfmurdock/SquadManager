package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"net/url"
	"time"

	"github.com/robertfmurdock/SquadManager/SquadManagerService/api"
	"github.com/stretchr/testify/assert"
)

type Tester struct {
	t       *testing.T
	Handler http.Handler
}

func New(t *testing.T, handler http.Handler) *Tester {
	return &Tester{t, handler}
}

func (tester *Tester) PerformRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	tester.Handler.ServeHTTP(recorder, request)
	return recorder
}

func (tester *Tester) DoRequest(method, urlStr string, body interface{}) Response {
	value, err := getPostBody(body)
	if err != nil {
		tester.t.Fatal(err)
	}
	bodyReader := bytes.NewReader(value)
	request := newRequest(tester.t, method, urlStr, bodyReader)
	return Response{tester, tester.PerformRequest(request)}
}

func (tester *Tester) GetSquadList(begin *time.Time, end *time.Time) Response {
	values := valuesWithDateRange(begin, end)
	squadUrl := tester.urlWithValues("/squad", values)
	return tester.DoRequest("GET", squadUrl.String(), nil)
}

func (tester *Tester) GetSquad(squadId api.SquadId, begin *time.Time, end *time.Time) Response {
	values := valuesWithDateRange(begin, end)
	return tester.GetSquadWithParameters(squadId, values)
}

func valuesWithDateRange(begin *time.Time, end *time.Time) *url.Values {
	values := &url.Values{}
	addTimeValue(begin, values, "begin")
	addTimeValue(end, values, "end")
	return values
}

func (tester *Tester) GetSquadWithParameters(squadId api.SquadId, values *url.Values) Response {
	squadUrl := tester.urlWithValues("/squad/"+squadId.String(), values)
	return tester.DoRequest("GET", squadUrl.String(), nil)
}
func (tester *Tester) urlWithValues(urlString string, values *url.Values) (*url.URL) {
	squadUrl, err := url.Parse(urlString)
	if err != nil {
		tester.t.Fatal(err)
	}
	squadUrl.RawQuery = values.Encode()
	return squadUrl
}

func addTimeValue(t *time.Time, values *url.Values, key string) {
	if t != nil {
		values.Add(key, api.FormatDate(t))
	}
}

func (tester *Tester) PostSquad() Response {
	return tester.DoRequest("POST", "/squad", "")
}

func (tester *Tester) PostSquadMember(squadId api.SquadId, member api.SquadMember) Response {
	return tester.DoRequest("POST", "/squad/"+squadId.String(), member)
}

func (tester *Tester) PerformPostSquad() api.SquadId {
	var newSquadId api.SquadId
	tester.PostSquad().
		CheckStatus(http.StatusAccepted).
		LoadJson(&newSquadId)
	return newSquadId
}

func (tester *Tester) PerformGetSquad(squadId api.SquadId, begin *time.Time, end *time.Time) api.Squad {
	squad := api.Squad{}
	tester.GetSquad(squadId, begin, end).
		CheckStatus(http.StatusOK).
		LoadJson(&squad)
	return squad
}

func (tester *Tester) PerformGetSquadList(begin *time.Time, end *time.Time) []api.Squad {
	var loadedJson []api.Squad
	tester.GetSquadList(begin, end).
		CheckStatus(http.StatusOK).
		LoadJson(&loadedJson)
	return loadedJson
}

func (tester *Tester) PerformPostSquadMember(squadId api.SquadId, squadMember api.SquadMember) {
	var newSquadMemberId api.SquadMemberId
	tester.PostSquadMember(squadId, squadMember).
		CheckStatus(http.StatusAccepted).
		LoadJson(&newSquadMemberId)

	assert.Equal(tester.t, squadMember.ID, newSquadMemberId)
}

type Response struct {
	Tester   *Tester
	Recorder *httptest.ResponseRecorder
}

func (response Response) CheckStatus(status int) Response {
	if !assert.Equal(response.Tester.t, status, response.Recorder.Code) {
		response.Tester.t.Fatal(response.Recorder.Body.String())
	}

	return response
}

func (response Response) LoadJson(loadLocation interface{}) Response {
	err := json.Unmarshal(response.Recorder.Body.Bytes(), loadLocation)
	if err != nil {
		response.Tester.t.Fatal(err, response.Recorder.Body.String())
	}
	return response
}

func newRequest(t *testing.T, method, urlStr string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func getPostBody(body interface{}) ([]byte, error) {
	if value, ok := body.(string); ok {
		return []byte(value), nil
	} else {
		return json.Marshal(body)
	}
}
