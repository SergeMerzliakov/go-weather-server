package service

import (
	"bytes"
	"github.com/SergeMerzliakov/go-weather-server/config"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const SYDNEY_DATA = "{\"coord\":{\"lon\":151.21,\"lat\":-33.87},\"weather\":[{\"id\":803,\"main\":\"Clouds\",\"description\":\"broken clouds\",\"icon\":\"04d\"}],\"base\":\"stations\",\"main\":{\"temp\":289.56,\"feels_like\":287.77,\"temp_min\":288.71,\"temp_max\":290.37,\"pressure\":1019,\"humidity\":77},\"visibility\":10000,\"wind\":{\"speed\":3.6,\"deg\":340},\"clouds\":{\"all\":76},\"dt\":1590882599,\"sys\":{\"type\":1,\"id\":9600,\"country\":\"AU\",\"sunrise\":1590871874,\"sunset\":1590908087},\"timezone\":36000,\"id\":2147714,\"name\":\"Sydney\",\"cod\":200}"
const SYDNEY = "Sydney"

func TestParsing(t *testing.T) {
	var tests = []struct {
		data          string
		requestCity   string
		errorExpected bool
		errorStatus   int
		errorResponse string
	}{
		{SYDNEY_DATA, SYDNEY, false, http.StatusOK, ""},
		{SYDNEY_DATA, SYDNEY, true, http.StatusInternalServerError, "Weather API error: 500"},
	}

	for _, tt := range tests {
		test := tt

		client := ClientMock{
			Cities:     []string{test.requestCity},
			HasError:   test.errorExpected,
			StatusCode: test.errorStatus,
		}
		ws := CreateWeatherService(&config.ServerConfiguration{}, &client)

		cities := []string{test.requestCity}
		reports, err := ws.GetCityWeather(cities)

		if !test.errorExpected {
			assert.Nil(t, err)
			assert.NotNil(t, reports)
			assert.Equal(t, 1, len(reports.Reports))
			assert.NotNil(t, reports.Reports[SYDNEY])

		} else {
			assert.NotNil(t, err)
			assert.Nil(t, reports)
			assert.Equal(t, test.errorResponse, err.Error())
		}
	}
}

type ClientMock struct {
	Cities     []string
	HasError   bool
	StatusCode int
}

func (mc *ClientMock) Do(req *http.Request) (*http.Response, error) {

	if mc.HasError {
		resp := http.Response{
			Status:     "Error",
			StatusCode: mc.StatusCode,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Error")),
		}
		return &resp, nil

	} else {
		resp := http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(SYDNEY_DATA)),
		}
		return &resp, nil
	}

}
