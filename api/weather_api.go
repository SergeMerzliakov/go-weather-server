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

package api

import (
	"encoding/json"
	"fmt"
	"github.com/SergeMerzliakov/go-weather-server/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type WeatherRequest struct {
	Cities []string
}

type SynchronousWeatherEndpoint struct {
	Service service.WeatherService
}

// CreateWeatherEndpoint is a factory method for creating endpoints and will also be used
// for testing by injecting mock weather services
func CreateWeatherEndpoint(svc service.WeatherService) *SynchronousWeatherEndpoint {
	return &SynchronousWeatherEndpoint{Service: svc}
}

// WeatherEndpoint is the top level handler for a URL path, and does all logging
func (h *SynchronousWeatherEndpoint) WeatherEndpoint(w http.ResponseWriter, r *http.Request) {
	wr := WeatherRequest{}
	reqBody, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Error("Error reading client HTTP request", readErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "Invalid Request format. Could not parse request body")
		return
	}
	// get list of cities requested
	jsonErr := json.Unmarshal(reqBody, &wr.Cities)
	if jsonErr != nil {
		log.Error("Error converting client weather request into json", jsonErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Invalid Request format. Check json format is an array of strings."))
		return
	}

	result, callErr := h.Service.GetCityWeather(wr.Cities)

	if callErr != nil {
		log.Error("Error invoking Weather API", callErr)
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("Weather API error"))
		return
	}

	data, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		log.Error("Error converting Weather API response into json", callErr)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Error marshalling JSON response"))
		return
	}

	log.Info("Successfully retreived weather request for client")
	_, _ = w.Write(data)
}
