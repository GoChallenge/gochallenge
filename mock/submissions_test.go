package mock_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/mock"
	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestSubmissionsMockRepo(t *testing.T) {
	ss := mock.NewSubmissions()
	s0 := &model.Submission{
		ID: "1",
	}
	require.NoError(t, ss.Add(s0), "adding s0 submission errored")

	s1 := &model.Submission{
		ID: "2",
	}
	require.NoError(t, ss.Add(s1), "adding s1 submission errored")

	// All should return all added submissions
	sx, err := ss.All()
	require.NoError(t, err, "all submissions returned an error")
	require.Equal(t, []*model.Submission{s0, s1}, sx,
		"all submissions not returned correctly")
}

func TestSubmissionsAllForChallenge(t *testing.T) {
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

	s0 := &model.Submission{
		ID:        "1",
		Challenge: c0,
	}
	require.NoError(t, ss.Add(s0), "adding s0 submission errored")

	c1 := &model.Challenge{
		ID: 2,
	}
	s1 := &model.Submission{
		ID:        "2",
		Challenge: c1,
	}
	require.NoError(t, ss.Add(s1), "adding s1 submission errored")

	// AllForChallenge should return all submissions for a given
	// challenge only
	sx, err = ss.AllForChallenge(c0)
	require.NoError(t, err, "all for challenge returned an error")
	require.Equal(t, []*model.Submission{s0}, sx,
		"all submissions not returned correctly")
}
