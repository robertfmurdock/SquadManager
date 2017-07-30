package api

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Squad struct {
	ID      string
	Members []SquadMember
}

type SquadMember struct {
	ID    string
	Range Range
	Email string
}

func NewSquadMember(email string, dateRange Range) SquadMember {
	return SquadMember{
		ID:    bson.NewObjectId().Hex(),
		Range: dateRange,
		Email: email,
	}
}

type Range struct {
	Begin time.Time
	End   time.Time
}

func (r Range) MarshalJSON() ([]byte, error) {
	rangeElement := struct {
		Begin string
		End   string
	}{
		FormatDate(&r.Begin),
		FormatDate(&r.End),
	}

	return json.Marshal(rangeElement)
}

func FormatDate(t *time.Time) string {
	return t.Format(time.RFC3339)
}

func ParseDate(date string) (*time.Time, error) {
	if len(date) == 0 {
		return nil, nil
	}
	parse, err := time.Parse(time.RFC3339, date)
	return &parse, err
}

func Date(year int, month time.Month, day int) *time.Time {
	location := time.Now().Location()
	date := time.Date(year, month, day, 0, 0, 0, 0, location)
	return &date
}

func FilterMembers(members []SquadMember, begin *time.Time, end *time.Time) []SquadMember {
	if begin == nil && end == nil {
		return members
	}
	var result []SquadMember

	for _, member := range members {
		if isAfterBeginning(begin, member) && isBeforeEnd(end, member) {
			result = append(result, member)
		}
	}

	return result
}
func isBeforeEnd(end *time.Time, member SquadMember) bool {
	return end == nil || member.Range.Begin.Before(*end)
}

func isAfterBeginning(begin *time.Time, member SquadMember) bool {
	return begin == nil || member.Range.End.After(*begin)
}
