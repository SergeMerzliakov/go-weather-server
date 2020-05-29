/*
 * Copyright 2020 Serge Merzliakov
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"encoding/json"
	"fmt"
	"github.com/SergeMerzliakov/go-weather-server/api"
	"github.com/SergeMerzliakov/go-weather-server/config"
	"github.com/SergeMerzliakov/go-weather-server/service"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	environment = "ENVIRONMENT"
	apiKey      = "API_KEY"
)

func main() {
	// initialize logging to STDOUT in json format for integration with log aggregators and processors like logstash
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	serverConfig, confErr := loadConfiguration()

	if confErr != nil {
		log.Error(errors.Wrap(confErr, "Cannot start server. Configuration error"))
		return
	}

	log.Debug("Starting weather server...")
	router := mux.NewRouter().StrictSlash(true)

	//TODO
	weatherApi := api.CreateWeatherEndpoint(service.CreateWeatherService(serverConfig, http.DefaultClient))

	// we use post so we can send request form with cities rather encode them as part of the request
	// parameters. This will support any number of requests
	router.HandleFunc("/weather", weatherApi.WeatherEndpoint).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverConfig.ServerPort), router))
}

// loadConfiguration
func loadConfiguration() (*config.ServerConfiguration, error) {
	env := os.Getenv(environment)

	var configFile = ""

	switch env {
	case "dev":
		configFile = "configFile/config.dev.json"
	case "test":
		configFile = "configFile/config.test.json"
	case "prod":
		configFile = "configFile/config.prod.json"
	default:
		return nil, errors.Errorf("Unknown environment set - '%s'", env)
	}

	log.Infof("loading configuration from '%s'", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}

	conf := config.ServerConfiguration{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	// get API key for openweathermap
	key := os.Getenv(apiKey)
	if len(key) == 0 {
		return nil, errors.Errorf("Bad API key set - '%s'", key)
	}

	conf.APIKey = key
	return &conf, nil
}
