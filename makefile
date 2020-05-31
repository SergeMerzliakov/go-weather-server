build:
	go build -o bin/main main.go

test:
	go test ./...

run:
	export ENVIRONMENT=dev
	echo "make sure API_KEY environment variable is set to value from https://openweathermap.org/"
	go run main.go