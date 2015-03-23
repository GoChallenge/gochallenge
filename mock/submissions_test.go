package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/gochallenge/gochallenge/model/spec"
	"github.com/stretchr/testify/require"
)

func TestSubmissionsMockRepoSpec(t *testing.T) {
	ss := mock.NewSubmissions()
	spec.MustBehaveLikeSubmissions(t, &ss)
}

func TestNoSubmissionsAllForChallenge(t *testing.T) {
	c0 := &model.Challenge{
		ID: 1,
	}
	ss := mock.NewSubmissions()

	// AllForChallenge should return empty array when there're
	// no submissions for a challenge
	sx, err := ss.AllForChallenge(c0)
	require.NoError(t, err, "all for challenge returned an error")
	require.Equal(t, []*model.Submission{}, sx,
		"empty submissions not returned correctly")
}
