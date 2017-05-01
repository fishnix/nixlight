package common

import (
	"os"
	"strings"
)

// Env parses the environment variables and returns them as a map
func Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}
