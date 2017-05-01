package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/fishnix/darksky"
	"github.com/fishnix/nixlight/common"
	"github.com/fishnix/nixlight/light"
	"github.com/fishnix/nixlight/nixlight"
)

var (
	config     = flag.String("config", "./nixlight.toml", "Configuration file.")
	version    = flag.Bool("version", false, "Display version information and exit.")
	buildstamp = "No Version Provided"
	githash    = "No Git Commit Provided"
)

// Version is the main version number
const Version = nixLight.Version

// VersionPrerelease is a prerelease marker
const VersionPrerelease = nixLight.VersionPrerelease

func main() {
	log.Println("[INFO] Starting NiXLight!")
	flag.Parse()

	if *version {
		vers()
	}

	// Read configuration file.
	configuration := common.ReadConfig(config)
	configuration.Print()

	// START CONSUMERS
	persister, _ := startPersister(configuration.DB)
	// END CONSUMERS

	// START FETCHER
	apiClient, err := buildAPIClient(configuration)
	if err != nil {
		log.Fatalln("[ERROR] unable to initialize the DarkSky API Client", err)
	}

	timer, err := time.ParseDuration(configuration.Timer)
	if err != nil {
		log.Fatalln("[ERROR] unable to parse configuration Timer", err)
	}

	log.Println("[INFO] Updating every", configuration.Timer)
	fetcher, _ := startFetcher(apiClient, timer)
	// END FETCHER

	// START FANOUT
	for f := range fetcher {
		persister <- f
	}
	// END COLLECTOR

	dummy := light.Dummy{Name: "Dummy Light"}
	dummy.Start()
	dummy.Brightness(123)

	if l, ok := interface{}(dummy).(light.Controller); ok {
		log.Printf("It looks like %+v satisfies the light.Controller interface!", l)
	}

	// sleep FOREVAH
	for {
		time.Sleep(time.Millisecond * 10000)
	}
	// close(quit)
}

func buildAPIClient(configuration common.Config) (darksky.APIClient, error) {
	apiClient := darksky.APIClient{
		Lat:     configuration.DarkSkyClient.Lat,
		Long:    configuration.DarkSkyClient.Long,
		Key:     configuration.DarkSkyClient.Key,
		Exclude: []string{"hourly", "minutely"},
	}

	return apiClient, nil
}

// startFetcher takes the darksky client and creates a goroutine that returns
// forecast data over a channel and opens a control channel for killing the goroutine
func startFetcher(client darksky.APIClient, timer time.Duration) (<-chan *darksky.Forecast, chan string) {
	controlChannel := make(chan string)
	forecastChannel := make(chan *darksky.Forecast)

	ticker := time.NewTicker(timer)
	go func() {
		for {
			select {
			case <-ticker.C:
				forecast, err := client.GetForecast()
				if err != nil {
					log.Println("[ERROR] unable to get forecast", err)
				}
				forecastChannel <- forecast
			case <-controlChannel:
				close(forecastChannel)
				ticker.Stop()
				return
			}
		}
	}()

	return forecastChannel, controlChannel
}

// startPersister takes the boltdb file and starts a goroutine that receives forecast
// data over a channel and opens a control channel for killing the goroutine
func startPersister(file string) (chan *darksky.Forecast, chan string) {
	initializeDB(file)
	controlChannel := make(chan string)
	forecastChannel := make(chan *darksky.Forecast)
	go func() {
		for {
			select {
			case f := <-forecastChannel:
				log.Println("[INFO] recieved forecast data on persister channel")
				log.Printf("Forecast: %+v", f)
				err := persistData(file, f)
				if err != nil {
					log.Println("[ERROR] unable to perist data:", err)
				}
			case <-controlChannel:
				close(forecastChannel)
				return
			}
		}
	}()

	return forecastChannel, controlChannel
}

func initializeDB(file string) {
	// open and/or create a boltdb
	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalln("[ERROR] unable to initialize the database", err)
	}
	defer db.Close()

	// initialize the forecast buckets
	initBucket(db, "DailyForecast")
	initBucket(db, "Currently")
}

func initBucket(db *bolt.DB, bucket string) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatalln("[ERROR] unable to initialize bucket:", err)
		}
		return nil
	})
}

func persistData(file string, forecast *darksky.Forecast) error {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		log.Fatal("[ERROR] Unable to open boltdb", err)
	}
	defer db.Close()

	cdata, err := json.Marshal(forecast.Currently)
	if err != nil {
		log.Println("[ERROR] Unable to marshall forecast into JSON:", err)
	}
	store(db, "Currently", []byte("now"), cdata)

	for _, d := range forecast.Daily.Data {
		tm := time.Unix(int64(d.Time), 0)
		date := []byte(tm.String())
		ddata, err := json.Marshal(d)
		if err != nil {
			log.Println("[ERROR] Unable to marshall forecast into JSON:", err)
		}
		store(db, "DailyForecast", date, ddata)
	}

	return nil
}

func store(db *bolt.DB, bucket string, key []byte, data []byte) error {
	db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucket))
		err = b.Put(key, data)
		return err
	})

	return nil
}

func vers() {
	fmt.Printf("Nix Light Version: %s%s\n", Version, VersionPrerelease)
	fmt.Println("Git Commit Hash:", githash)
	fmt.Println("UTC Build Time:", buildstamp)
	os.Exit(0)
}
