package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model/spec"
)

func TestUsersSpec(t *testing.T) {
	us := mock.NewUsers()
	spec.MustBehaveLikeUsers(t, &us)
}
