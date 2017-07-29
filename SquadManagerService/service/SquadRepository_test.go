package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestRepositoryFactoryWillCopySession(t *testing.T) {
	factory := SquadRepositoryFactory{
		Config: Configuration{
			Host:         "127.0.0.1",
			DatabaseName: "SquadRepositoryTest",
			DbTimeout:    time.Second / 100,
		}, parentSession: nil,
	}

	repository1, err1 := factory.Repository()
	repository2, err2 := factory.Repository()

	if err1 != nil {
		t.Fatal(err1)
	}

	if err2 != nil {
		t.Fatal(err2)
	}

	defer repository1.Close()
	defer repository2.Close()

	assert.False(t, repository1.session == factory.parentSession)
	assert.False(t, repository1.session == repository2.session)
}
