package api

import (
	"time"
	"encoding/json"
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

type Range struct {
	Begin time.Time
	End   time.Time
}

func (r Range) MarshalJSON() ([]byte, error) {
	rangeElement := struct {
		Begin string
		End   string
	}{
		r.Begin.Format(time.RFC3339),
		r.End.Format(time.RFC3339),
	}

	return json.Marshal(rangeElement)
}
