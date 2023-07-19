package config

import (
	"sync"

	"github.com/stevenroose/gonfig"
)

type Config struct {
	DBUsername string `default:"root"`
	//DBPassword string `default:"hz990129"`
	DBPassword string `default:"123456"`
	DBAddress  string `default:"127.0.0.1:3306"`
	DBName     string `default:"uvm"`
	PORT       string `default:"8000"`

	// Use 'id' to change the name of a flag.
	ConfigFile string `id:"config" short:"c"`
}

var once sync.Once
var config Config

func GetConfig() *Config {
	err := gonfig.Load(&config, gonfig.Conf{
		ConfigFileVariable:  "config", // enables passing --configfile myfile.conf
		FileDefaultFilename: "uvm.json",
		// The default decoder will try TOML, YAML and JSON.
		FileDecoder: gonfig.DecoderJSON,

		EnvPrefix: "UVM_",
	})
	if err != nil {
		panic(err)
	}
	return &config
}
