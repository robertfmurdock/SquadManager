package api

type Squad struct {
	ID string
	Members []SquadMember
}

type SquadMember struct {
	ID string
	Email string
}
