package persist

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fishnix/darksky"
)

// BoltDbPersister defines the struct for persisting data in BoltDB
type BoltDbPersister struct {
	File    string
	Options bolt.Options
}

// Initialize prepares the BoltDB for persisting
func (db BoltDbPersister) Initialize() error {
	// open and/or create a boltdb
	boltdb, err := bolt.Open(db.File, 0600, &db.Options)
	if err != nil {
		log.Println("[ERROR] unable to initialize the database", err)
		return err
	}
	defer boltdb.Close()

	// initialize the forecast buckets
	initBucket(boltdb, "DailyForecast")
	initBucket(boltdb, "Currently")
	initBucket(boltdb, "Configuration")

	return nil
}

// initBucket initializes a bucket in BoltDB
func initBucket(db *bolt.DB, bucket string) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatalln("[ERROR] unable to initialize bucket:", err)
		}
		return nil
	})
}

// Persist is the interface method of Persister to write data to a bolt db
func (db BoltDbPersister) Persist(forecast *darksky.Forecast) error {
	boltdb, err := bolt.Open(db.File, 0600, &db.Options)
	if err != nil {
		log.Fatal("[ERROR] Unable to open boltdb", err)
	}
	defer boltdb.Close()

	cdata, err := json.Marshal(forecast.Currently)
	if err != nil {
		log.Println("[ERROR] Unable to marshall forecast into JSON:", err)
	}
	store(boltdb, "Currently", []byte("now"), cdata)

	for _, d := range forecast.Daily.Data {
		tm := time.Unix(int64(d.Time), 0)
		date := []byte(tm.String())
		ddata, err := json.Marshal(d)
		if err != nil {
			log.Println("[ERROR] Unable to marshall forecast into JSON:", err)
		}
		store(boltdb, "DailyForecast", date, ddata)
	}

	return nil
}

// store does the boltdb Update
func store(db *bolt.DB, bucket string, key []byte, data []byte) error {
	db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucket))
		err = b.Put(key, data)
		return err
	})

	return nil
}
