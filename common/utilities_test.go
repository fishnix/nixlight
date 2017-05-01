package common_test

import (
	"os"
	"testing"

	"github.com/fishnix/nixlight/common"
)

func TestEnv(t *testing.T) {
	os.Clearenv()

	envMap := map[string]string{
		"foo": "bar",
		"fiz": "baz",
		"biz": "buz",
	}

	for k, v := range envMap {
		os.Setenv(k, v)
	}

	env := common.Env()
	for k, v := range env {
		if env[k] != v {
			t.Error("Expected env[k] to be ", v, "got:", env[k])
		}
	}
}
