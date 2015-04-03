package boltdb

import (
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gochallenge/gochallenge/model"
)

type boltf func(tx *bolt.Tx) error

// Open Bolt database at the given file name
func Open(file string) (*bolt.DB, error) {
	opt := &bolt.Options{
		Timeout: 1 * time.Second,
	}
	return bolt.Open(file, 0600, opt)
}

// chain executes a given chain of database operations within a single
// transaction, which is provided as the first argument
func chain(tx func(func(tx *bolt.Tx) error) error, ops ...boltf) error {
	return tx(func(tx *bolt.Tx) error {
		var err error
		for _, op := range ops {
			if err = op(tx); err != nil {
				break
			}
		}
		return err
	})
}

// initialises bolt bucket for challenges
func initBucket(bkt []byte) boltf {
	return func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bkt)
		return err
	}
}

// returns a value stored in the bucket under the given integer key
func itoBytes(bkt []byte, id int, b **[]byte) boltf {
	return func(tx *bolt.Tx) error {
		k := strconv.Itoa(id)
		v := tx.Bucket(bkt).Get([]byte(k))

		// bolt returns an empty result for unknown key lookup,
		// return ErrNotFound in this case
		if v == nil {
			return model.ErrNotFound
		}
		*b = &v
		return nil
	}
}
