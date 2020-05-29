# go-weather-server
Simple weather API in Go

# Configuration

There are two kinds of configuration - non-sensitive and sensitive

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


### Sensitive or dynamic configuration
This is information which cannot be stored in files in source control (API keys, passwords) and are passed as environment variables

| KEY | DESCRIPTION | SAMPLE VALUE|
|-----|-----| ----|
| ENVIRONMENT | execution environment| dev, test, prod|
| API_KEY | key for openweathermap.org service|0f92045f1aa7432099325ce2f3b022e1ef0 |


## Things To Improve

* Needs Unit and Integration Tests
* used dep for dependency management and not go modules
* environment specific configuration, possibly dynamically determined