package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/gochallenge/gochallenge/model/spec"
)

func TestChallengeMockRepo(t *testing.T) {
	cs := mock.NewChallenges()
	cur := model.Challenge{
		ID: mock.CurrentID,
	}
	spec.MustBehaveLikeChallenges(t, &cs, &cur)
}
