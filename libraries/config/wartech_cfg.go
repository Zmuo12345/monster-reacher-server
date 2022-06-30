package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const PATH_PROJECT = "D:/WTProject/monster-reacher/server"

type wartechConfig struct {
	Services map[string]wartechConfigServices `json:"services"`
}

type wartechConfigServices struct {
	Hosts []string `json:"hosts"`
	Ports []int    `json:"ports"`
}

var cacheWartechConfig *wartechConfig = initWartechConfig()

func initWartechConfig() *wartechConfig {
	cfg := &wartechConfig{}
	jsonFile, err := os.Open(PATH_PROJECT + "/config/wartech_config.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &cfg)
	return cfg
}

func WartechConfig() wartechConfig {
	return *cacheWartechConfig
}
