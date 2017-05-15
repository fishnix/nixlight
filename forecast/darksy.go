package forecast

import (
	"github.com/fishnix/darksky"
)

// DarkSkyForecast defines the struct for getting forecast data from the Darksky API
type DarkSkyForecast struct {
	Lat     string
	Long    string
	Key     string
	Exclude []string
}

// Initialize prepares the BoltDB for persisting
func (d DarkSkyForecast) Initialize() error {
	return nil
}

// Fetch gets the forecast from the darkski api
// TODO: parse this into a common format, should return nixlight.Forecasts
func (d DarkSkyForecast) Fetch() (*darksky.Forecast, error) {
	client := d.buildAPIClient()
	return client.GetForecast()
}

func (d DarkSkyForecast) buildAPIClient() darksky.APIClient {
	apiClient := darksky.APIClient{
		Lat:     d.Lat,
		Long:    d.Long,
		Key:     d.Key,
		Exclude: d.Exclude,
	}

	return apiClient
}
