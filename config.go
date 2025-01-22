package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type MyForwardProxyConfig struct {
	ExpirationPeriod int `toml:"expiration_period"`
}

type Config struct {
	MyForwardProxy MyForwardProxyConfig `toml:"my-forward-proxy"`
}

func loadConfig() Config {
	config := Config{}

	file, err := os.Open("config.toml")
	if err != nil {
		log.Fatalf("no config file found, you can generate one with ./bin/my_forward_proxy --init or make config")
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	_, err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding TOML file: %v", err)
	}

	return config
}

func initializeConfig() {
	defaultConfig := MyForwardProxyConfig{
		ExpirationPeriod: 3600,
	}

	configMap := Config{
		MyForwardProxy: defaultConfig,
	}

	file, err := os.Create("config.toml")
	if err != nil {
		fmt.Printf("Error creating config.toml: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(configMap); err != nil {
		fmt.Printf("Error writing to config.toml: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config created")
}
