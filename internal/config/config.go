package config


//error handling is must important in go

import (
	"os"
	"log"
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address string `yaml:"address"`
}


//env-default:"production"


type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

//MUST() functions m error through(return) nhi karna chye if function nhi chala toh vnhi se fatal return kr do and stop the program execution
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
		flags := flag.String("config", "", "Path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not provided. Use --config flag or set CONFIG_PATH environment variable.")
		}
	}


	//_, err this id short form ki info not need only checkerror
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	}

	return &cfg
}