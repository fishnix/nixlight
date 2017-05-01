package common_test

import (
	"testing"

	"github.com/fishnix/nixlight/common"
)

func TestReadConfig(t *testing.T) {
	title := "nixlight test config"
	clientKey := "changeme"
	clientLat := "51.894444"
	clientLong := "1.482500"

	configFile := "../example/test.toml"
	configuration := common.ReadConfig(&configFile)

	if configuration.Title != title {
		t.Error("Expected title to be:", title, " got:", configuration.Title)
	}

	if configuration.DarkSkyClient.Key != clientKey {
		t.Error("Expected client.Key to be:", clientKey, " got:", configuration.DarkSkyClient.Key)
	}

	if configuration.DarkSkyClient.Lat != clientLat {
		t.Error("Expected client.Location to be:", clientLat, " got:", configuration.DarkSkyClient.Lat)
	}

	if configuration.DarkSkyClient.Long != clientLong {
		t.Error("Expected client.Location to be:", clientLong, " got:", configuration.DarkSkyClient.Long)
	}
}

func TestConfig_Print(t *testing.T) {}

// sliceEq compares two slices for equality and returns a boolean
func slicesEq(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if (a == nil || b == nil) || (len(a) != len(b)) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
