package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//const PATH_PROJECT = "D:/WTProject/monster-reacher/server"
const PATH_PROJECT = "/app/monster-reacher-server/"

type wartechConfig struct {
	Services  map[string]wartechConfigServices `json:"services"`
	Databases map[string]wartechConfigDatabase `json:"databases"`
}

type wartechConfigServices struct {
	Hosts []string `json:"hosts"`
	Ports []int    `json:"ports"`
}

type wartechConfigDatabase struct {
	Host string `json:"host"`
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
