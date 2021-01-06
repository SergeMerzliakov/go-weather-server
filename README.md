# go-weather-server
A simple weather API in Go. It is a proxy for a free weather service (openweathermap).

It was written in 3 hours as a code test for a developer role, I did not end up applying for :-(

# Configuration

There are two kinds of configuration - non-sensitive and sensitive.

### Non-Sensitive Configuration
These are stored as Json files in the directory and include things which COULD be stored in source control or something similar
 
    configData/
    

There should one for each environment.

#### Sample Configuration

    {
      "apiUrl": "http://api.openweathermap.org/data/2.5/weather?q=",
      "apiKeyParam": "appId",
      "serverPort" : 8080
    }


### Prerequisites

Install:
* Make utility 
* Go version 1.14 or better
* dep tool


### Running the Service

* Build the service with:
    

        make build
    

* Run the unit tests

     
        make test
        
* Get an API key (sample below) from openweathermap.org and add it to environment:
 

        export APL_KEY=0f92045f1aa7432099325ce2f3b022e1ef0
    

* Set the ENVIRONMENT to value (makefile run target defaults this to "dev")

* Run the service locally:


        make run
    


### Sensitive or dynamic configuration
This is information which cannot be stored in files in source control (API keys, passwords) and are passed as environment variables

| KEY | DESCRIPTION | SAMPLE VALUE|
|-----|-----| ----|
| ENVIRONMENT | execution environment| dev, test, prod|
| API_KEY | key for openweathermap.org service|0f92045f1aa7432099325ce2f3b022e1ef0 |



## Request

If running on your machine, call the POST endpoint /weather with:

    http://localhost:8080/weather


### Request Body
Array of strings, each of a city (case insenstive). For example:


    [
    "sydney",
    "adelaide",
    "Gotham City"
    ]



### Response

The weather details for each city are returned


    {
       "reports":{
          "Gotham City":{
             "description":"not found"
          },
          "adelaide":{
             "description":"clear sky",
             "temperature":16,
             "humidity":41
          },
          "sydney":{
             "description":"few clouds",
             "temperature":15,
             "humidity":82
          }
       }
    }



## Things To Improve

* Needs more unit and integration tests
* Used dep for dependency management, and not modules
* Environment specific configuration, possibly dynamically determined
* API spec (such as OpenAPI or Swagger)
* Sample Postman collections would be nice
* Dockerize
