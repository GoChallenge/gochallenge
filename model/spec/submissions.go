package spec

import (
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

// MustBehaveLikeSubmissions tests behaviour of the submission repo
// implementation, to make sure it conforms to the spec
func MustBehaveLikeSubmissions(t *testing.T, ss model.Submissions) {
	c1 := model.Challenge{
		ID: 123,
	}

	s1 := model.Submission{
		ID:        "01-c001c0ffee",
		Type:      model.LvlBonus,
		Created:   time.Now(),
		Challenge: &c1,
	}

	// Find should error before submission is added
	_, err := ss.Find(s1.ID)
	require.Equal(t, model.ErrNotFound, err)

	// Add should succeed
	err = ss.Add(&s1)
	require.NoError(t, err)

	// And now we should be able to find the same record by its ID
	sx, err := ss.Find(s1.ID)
	require.NoError(t, err)
	require.Equal(t, s1, *sx)

	// Let's add another submission, for another challenge this time
	c2 := model.Challenge{
		ID: 987,
	}
	s2 := model.Submission{
		ID:        "02-badc0ffee",
		Type:      model.LvlNormal,
		Created:   time.Now(),
		Challenge: &c2,
	}
	err = ss.Add(&s2)
	require.NoError(t, err)

	// All should return both submissions
	sxs, err := ss.All()
	require.NoError(t, err)
	require.Equal(t, 2, len(sxs))
	sx1 := *sxs[0]
	sx2 := *sxs[1]
	require.True(t, (sx1 == s1 && sx2 == s2) || (sx1 == s2 && sx2 == s1))

	// But AllForChallenge should return one submission only
	sxs, err = ss.AllForChallenge(&c1)
	require.NoError(t, err)
	require.Equal(t, 1, len(sxs))
	require.True(t, *sxs[0] == s1)
}
