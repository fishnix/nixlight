package common

import (
	"bufio"
	"bytes"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

// Config is the exported nixlight configuration data
type Config struct {
	Title         string              `toml:"title"`
	Timer         string              `toml:"timer"`
	DB            string              `toml:"db"`
	DarkSkyClient DarkSkyClientConfig `toml:"DarkSky"`
}

// DarkSkyClientConfig describes the configuration for the dark sky client
type DarkSkyClientConfig struct {
	Key  string `toml:"apiKey"`
	Lat  string
	Long string
}

// ReadConfig takes a (toml) file, reads it from disk, and unmarshalls it.
func ReadConfig(file *string) Config {
	log.Println("Reading config from", *file)

	tomlBlob, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalln("Unable to read configuration file", err)
	}

	tmpl, err := template.New("config").Parse(string(tomlBlob))
	if err != nil {
		log.Fatalln("Unable to parse configuration file", err)
	}

	var b bytes.Buffer
	parsedTemplate := bufio.NewWriter(&b)
	err = tmpl.Execute(parsedTemplate, Env())
	if err != nil {
		log.Fatalln("Unable to substitute template values!", err)
	}
	parsedTemplate.Flush()

	var config Config
	err = toml.Unmarshal(b.Bytes(), &config)

	if err != nil {
		log.Println("Parsed configuration file:", string(b.Bytes()))
		log.Fatalln("Unable to unmarshall TOML from configuration file", err)
	}

	return config
}

// Print prints a Config.
func (c *Config) Print() {
	log.Printf("%+v", c)
}
