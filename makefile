export ENVIRONMENT=dev

build:
	go build -o bin/main main.go

test:
	go test ./...

run:
	@echo
	@echo "NOTE!! Make sure API_KEY environment variable is set to value from https://openweathermap.org/"
	@echo
	go run main.go