package boltdb_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gochallenge/gochallenge/boltdb"
	"github.com/gochallenge/gochallenge/model/spec"
	"github.com/stretchr/testify/require"
)

func TestUsersBoltRepo(t *testing.T) {
	f, err := ioutil.TempFile("", "gctestboltdb")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	db, err := boltdb.Open(f.Name())
	require.NoError(t, err)
	us, err := boltdb.NewUsers(db)
	require.NoError(t, err)

	spec.MustBehaveLikeUsers(t, &us)
}
