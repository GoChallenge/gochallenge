package boltdb_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/gochallenge/gochallenge/boltdb"
	"github.com/gochallenge/gochallenge/model"
	"github.com/gochallenge/gochallenge/model/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChallengeBoltRepo(t *testing.T) {
	f, err := ioutil.TempFile("", "gctestboltdb")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	db, err := boltdb.Open(f.Name())
	require.NoError(t, err)
	cs, err := boltdb.NewChallenges(db)
	require.NoError(t, err)

	cur := model.Challenge{
		ID:    model.ChallengeID(rand.Intn(100) + 1e3),
		Name:  "Currently Running Challenge",
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
	}
	assert.True(t, cur.Current(), "current challenge must be current")

	spec.MustBehaveLikeChallenges(t, &cs, &cur)
}
