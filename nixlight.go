package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fishnix/nixlight/common"
	"github.com/fishnix/nixlight/forecast"
	"github.com/fishnix/nixlight/light"
	"github.com/fishnix/nixlight/nixlight"
	"github.com/fishnix/nixlight/persist"
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

	// Start Persisters and Consumers of forecast data
	db := persist.BoltDbPersister{File: configuration.DB}
	persister, _ := persist.StartPersister(db)

	ds := forecast.DarkSkyForecast{
		Lat:     configuration.DarkSkyClient.Lat,
		Long:    configuration.DarkSkyClient.Long,
		Key:     configuration.DarkSkyClient.Key,
		Exclude: []string{"hourly", "minutely"},
	}

	timer, err := time.ParseDuration(configuration.Timer)
	if err != nil {
		log.Fatalln("[ERROR] unable to parse configuration Timer", err)
	}
	log.Println("[INFO] Updating every", configuration.Timer)

	// Start the forecast fetcher
	fetcher, _ := forecast.StartFetcher(ds, timer)

	// Fanout fetcher to all consumers/persisters
	for f := range fetcher {
		persister <- f
	}

	dummy := light.Dummy{Name: "Dummy Light"}
	dummy.Start()
	dummy.SetBrightness(123)

	if l, ok := interface{}(dummy).(light.Controller); ok {
		log.Printf("It looks like %+v satisfies the light.Controller interface!", l)
	}

	// sleep FOREVAH
	for {
		time.Sleep(time.Millisecond * 10000)
	}
	// close(quit)
}

func vers() {
	fmt.Printf("Nix Light Version: %s%s\n", Version, VersionPrerelease)
	fmt.Println("Git Commit Hash:", githash)
	fmt.Println("UTC Build Time:", buildstamp)
	os.Exit(0)
}
