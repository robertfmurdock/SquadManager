package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryFactoryWillCopySession(t *testing.T) {

	factory := SquadRepositoryFactory{
		Config: Configuration{
			Host:         "0.0.0.0",
			DatabaseName: "SquadRepositoryTest",
			DbTimeout:    time.Second / 100,
		}, parentSession: nil,
	}
	defer factory.Close()

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

	assert.NotNil(t, factory.parentSession)
	assert.False(t, repository1.session == factory.parentSession)
	assert.False(t, repository1.session == repository2.session)
}
