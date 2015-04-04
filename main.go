package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gochallenge/gochallenge/api"
	"github.com/gochallenge/gochallenge/boltdb"
	"github.com/gochallenge/gochallenge/github"
	"github.com/gochallenge/gochallenge/mock"
)

const boltdbMode = 0600
const boltdbTimeout = 5 * time.Second

func main() {
	var (
		dbpath string
		port   int
	)
	rand.Seed(time.Now().UTC().UnixNano())

	flag.StringVar(&dbpath, "db", os.TempDir()+"gochal.db",
		"full path to the location of database file")
	flag.IntVar(&port, "port", 8081, "port to listen on")
	flag.Parse()

	// if database file doesn't exist - we should seed it with
	// initial data, so let's save the file status before we opened it
	_, dbst := os.Stat(dbpath)

	db := open(dbpath)
	defer db.Close()

	cfg := config(db)
	fmt.Printf("dbst: %+v\n", dbst)
	if os.IsNotExist(dbst) {
		fmt.Println("seeding the database")
		seedChallenges(cfg.Challenges)
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port),
		api.New(cfg)))
}

// open bolt database at the given path
func open(path string) *bolt.DB {
	db, err := bolt.Open(path, boltdbMode, &bolt.Options{
		Timeout: boltdbTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// create dependency configuration for the service
func config(db *bolt.DB) api.Config {
	cs, err := boltdb.NewChallenges(db)
	if err != nil {
		log.Fatal(err)
	}
	us, err := boltdb.NewUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	ss := mock.NewSubmissions()
	gh := github.NewClient()

	return api.Config{
		Challenges:  &cs,
		Submissions: &ss,
		Users:       &us,
		Github:      &gh,
	}
}
