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
		ID: 1,
	}
	require.NoError(t, ss.Add(s0), "adding s0 submission errored")

	s1 := &model.Submission{
		ID: 2,
	}
	require.NoError(t, ss.Add(s1), "adding s1 submission errored")

	// All should return all added submissions
	sx, err := ss.All()
	require.NoError(t, err, "all submissions returned an error")
	require.Equal(t, []*model.Submission{s0, s1}, sx,
		"all submissions not returned correctly")
}
