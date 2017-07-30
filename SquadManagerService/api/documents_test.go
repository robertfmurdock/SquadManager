package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterMembersByDateRange_NoFilter(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 8, 10),
			}),
		NewSquadMember("chip@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 11, 10),
			}),
	}

	filteredMembers := FilterMembers(members, nil, nil)

	assert.Equal(t, filteredMembers, members)
}

func TestFilterMembersByDateRange_OnlyBegin(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 8, 10),
			}),
		NewSquadMember("chip@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 11, 10),
			}),
	}

	filteredMembers := FilterMembers(members, Date(2017, 9, 0), nil)

	assert.Equal(t, filteredMembers, []SquadMember{members[1]})
}

func TestFilterMembersByDateRange_OnlyBegin_RightOnStart(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 8, 10),
			}),
		NewSquadMember("chip@fake.com",
			Range{
				Begin: *Date(2017, 9, 0),
				End:   *Date(2017, 11, 10),
			}),
	}

	filteredMembers := FilterMembers(members, Date(2017, 9, 0), nil)

	assert.Equal(t, filteredMembers, []SquadMember{members[1]})
}

func TestFilterMembersByDateRange_OnlyEnd(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 8, 10),
			}),
		NewSquadMember("chip@fake.com",
			Range{
				Begin: *Date(2017, 9, 1),
				End:   *Date(2017, 11, 10),
			}),
	}

	filteredMembers := FilterMembers(members, nil, Date(2017, 9, 0))

	assert.Equal(t, filteredMembers, []SquadMember{members[0]})
}

func TestFilterMembersByDateRange_IgnoresMembersTooEarly(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 7, 30),
				End:   *Date(2017, 8, 10),
			}),
		NewSquadMember("chip@fake.com",
			Range{
				Begin: *Date(2017, 9, 1),
				End:   *Date(2017, 11, 10),
			}),
	}

	filteredMembers := FilterMembers(members, Date(2017, 8, 20), Date(2017, 9, 0))

	assert.Equal(t, 0, len(filteredMembers))
}

func TestFilterMembersByDateRange_IncludesMembersGreaterThanGivenRange(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2016, 7, 30),
				End:   *Date(2020, 8, 10),
			}),
	}

	filteredMembers := FilterMembers(members, Date(2017, 8, 20), Date(2017, 9, 0))

	assert.Equal(t, members, filteredMembers)
}

func TestFilterMembersByDateRange_IncludesMembersSmallerThanGivenRange(t *testing.T) {
	members := []SquadMember{
		NewSquadMember("dale@fake.com",
			Range{
				Begin: *Date(2017, 8, 21),
				End:   *Date(2017, 8, 22),
			}),
	}

	filteredMembers := FilterMembers(members, Date(2017, 8, 20), Date(2017, 9, 0))

	assert.Equal(t, members, filteredMembers)
}
