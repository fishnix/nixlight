package persist

import (
	"log"

	"github.com/fishnix/darksky"
)

// Persister is an interface to persisting data from the darksky api
type Persister interface {
	Initialize() error
	Persist(*darksky.Forecast) error
}

// StartPersister takes a persister and starts a goroutine that receives forecast
// data over a channel and opens a control channel for killing the goroutine
func StartPersister(p Persister) (chan *darksky.Forecast, chan string) {
	err := p.Initialize()
	if err != nil {
		log.Fatalln("Error initializing persister", err)
	}

	controlChannel := make(chan string)
	forecastChannel := make(chan *darksky.Forecast)
	go func() {
		for {
			select {
			case f := <-forecastChannel:
				log.Println("[INFO] recieved forecast data on persister channel")
				log.Printf("%+v", f)
				err := p.Persist(f)
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
