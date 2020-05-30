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
package service

import (
	"encoding/json"
	"fmt"
	"github.com/SergeMerzliakov/go-weather-server/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type WeatherService interface {
	GetCityWeather(cities []string) (*WeatherReports, error)
}

type SynchronousWeatherService struct {
	config *config.ServerConfiguration
	Client *http.Client
}

// CreateWeatherService is a factory method to create a weather service and will also be used
// for testing by injecting mock configuration and http client
func CreateWeatherService(config *config.ServerConfiguration, client *http.Client) WeatherService {
	ws := &SynchronousWeatherService{config: config, Client: client}
	return ws
}

// GetCityWeather returns weather for a collection of cities
func (wc *SynchronousWeatherService) GetCityWeather(cities []string) (*WeatherReports, error) {
	result := WeatherReports{
		Reports: make(map[string]*CityReport),
	}

	// loop and call repeatedly for now. There may be ways to make a single call to do this
	for _, city := range cities {
		report, err := wc.getSingleCityWeather(city)
		if err != nil {
			return nil, err
		}
		result.Reports[city] = report
	}
	return &result, nil
}

// getSingleCityWeather returns weather for a single city. If city is not found (404), an empty "not found" Report is return
func (wc *SynchronousWeatherService) getSingleCityWeather(city string) (*CityReport, error) {
	log.Debugf("Fetching weather for city '%s'", city)

	fullUrl := wc.buildUrl(city)

	log.Info(fullUrl)
	req, reqRrr := http.NewRequest("GET", fullUrl, nil)
	if reqRrr != nil {
		return nil, reqRrr
	}

	resp, callErr := wc.Client.Do(req)
	if callErr != nil {
		return nil, callErr
	}

	if resp.StatusCode == http.StatusOK {
		return createReport(resp)
	} else if resp.StatusCode == http.StatusNotFound {
		log.Warnf("City '%s' not found. No weather data returned.", city)
		return &CityReport{
			Description: "not found",
		}, nil
	}

	return nil, errors.Errorf("Weather API error: %d", resp.StatusCode)
}

// buildUrl build url with all values including city name (case insensitive) and API key
// this is very simplistic and just appends strings together instead of using Url struct
func (wc *SynchronousWeatherService) buildUrl(city string) string {
	return fmt.Sprintf("%s%s&%s=%s", wc.config.API, city, wc.config.APIKeyParam, wc.config.APIKey)
}

// createReport maps the successful response for a single city to a CityReport
func createReport(resp *http.Response) (*CityReport, error) {
	// deserialize into map, so we can handle any data format and then
	// pick out the items we want
	bytes, ioErr := ioutil.ReadAll(resp.Body)
	if ioErr != nil {
		return nil, ioErr
	}

	var result map[string]interface{}
	unmErr := json.Unmarshal(bytes, &result)
	if unmErr != nil {
		return nil, unmErr
	}

	// build weather report from certain elements of API response
	weatherData := result["weather"].([]interface{})[0].(map[string]interface{})
	mainData := result["main"].(map[string]interface{})

	report := &CityReport{
		Description: weatherData["description"].(string),
		Temperature: formatTemperature(mainData["temp"]), // convert from absolute temperature
		Humidity:    mainData["humidity"].(float64),
	}

	return report, nil
}

// formatTemperature converts temperature from absolute temperature to degrees celsius
// error handling needs to be better
func formatTemperature(temp interface{}) float64 {
	num := temp.(float64)
	val, err := strconv.ParseFloat(fmt.Sprintf("%.0f", num-273.0), 64)
	if err != nil {
		log.Errorf("Error converting temperature string [%s]: %s", num, err)
		return -1.0 // todo need to do something better than this!
	}
	return val
}
