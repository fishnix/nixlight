package forecast

import (
	"log"
	"time"

	"github.com/fishnix/darksky"
)

// Fetcher is an interface for Fetching forecast data
type Fetcher interface {
	Initialize() error
	Fetch() (*darksky.Forecast, error)
}

// StartFetcher takes a Fetcher and a duration and creates a goroutine that returns
// forecast data over a channel and a control channel for killing the goroutine
func StartFetcher(f Fetcher, timer time.Duration) (<-chan *darksky.Forecast, chan string) {
	controlChannel := make(chan string)
	forecastChannel := make(chan *darksky.Forecast)

	ticker := time.NewTicker(timer)
	go func() {
		for {
			select {
			case <-ticker.C:
				forecast, err := f.Fetch()
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
