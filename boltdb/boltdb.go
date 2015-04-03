package boltdb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
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

// store record in the bucket under the given key
func store(bkt []byte, k interface{}, u interface{}) boltf {
	return func(tx *bolt.Tx) error {
		var b bytes.Buffer

		kb, err := vtoKey(k)
		if err != nil {
			return err
		}
		if err := gob.NewEncoder(&b).Encode(u); err != nil {
			return err
		}
		return tx.Bucket(bkt).Put(kb, b.Bytes())
	}
}

// retrieve a record stored in the bucket under the given key
func load(bkt []byte, k interface{}, u interface{}) boltf {
	return func(tx *bolt.Tx) error {
		var b *[]byte

		kb, err := vtoKey(k)
		if err != nil {
			return err
		}
		if err := atoBytes(bkt, kb, &b)(tx); err != nil {
			return err
		}
		return decode(b, u)
	}
}

// iterates through all objects in the bucket, returning the first
// one matching given predicate function
func first(bkt []byte, f func(interface{}) bool, x interface{}) boltf {
	return func(tx *bolt.Tx) error {
		var err error
		bk := tx.Bucket(bkt).Cursor()

		for k, v := bk.First(); k != nil && err == nil; k, v = bk.Next() {
			if err = decode(&v, x); err == nil && f(x) {
				// the matching record is found, stop here
				return nil
			}
		}
		// no matching record was found, if there're no errors either -
		// return ErrNotFound
		if err == nil {
			err = model.ErrNotFound
		}
		return err
	}
}

// returns a value stored in the bucket under the given string key
func atoBytes(bkt []byte, k []byte, b **[]byte) boltf {
	return func(tx *bolt.Tx) error {
		v := tx.Bucket(bkt).Get(k)

		// bolt returns an empty result for unknown key lookup,
		// return ErrNotFound in this case
		if v == nil {
			return model.ErrNotFound
		}
		*b = &v
		return nil
	}
}

// decodes a record from the given byte slice
func decode(b *[]byte, u interface{}) error {
	return gob.NewDecoder(bytes.NewReader(*b)).Decode(u)
}

// converts given value into its key representation.
// Key are stored as big-endian-encoded binary values, to allow
// Bolt to sort them automatically
func vtoKey(id interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, id)
	return b.Bytes(), err
}

// converts a binary into into its value representation
func keytoV(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	return binary.Read(r, binary.BigEndian, v)
}

// find the maximum key value in the bucket
func maxKey(tx *bolt.Tx, bkt []byte, id interface{}) error {
	// Bolt stores its key in an ordered fashion, which means
	// (as we store our keys as big-endian byte arrays, too)
	// we can simply grab the latest key and use its value
	bk := tx.Bucket(bkt)
	k, _ := bk.Cursor().Last()
	return keytoV(k, id)
}
