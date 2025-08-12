package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port         int    `json:"port"`
	Env          string `json:"env"`
	AppName      string `json:"appName"`
}

func NewConfig() *Config {
	conf := &Config{}
	err := conf.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return conf
}

func (conf *Config) Load() error {
	dir, _ := os.Getwd()

	dat, err := os.ReadFile(dir + "/config/config.json")

	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, conf)

	return err
}
