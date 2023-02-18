.DEFAULT_GOAL := start
APP_ANTENNA_LISTENER_IP ?= "0.0.0.0:4002"
API_SERVER_LISTENER_IP ?= "0.0.0.0:8080"

start:
	go build -o ./binaries/chicha ./chicha.go
	./binaries/chicha --collector $(APP_ANTENNA_LISTENER_IP) --web $(API_SERVER_LISTENER_IP)

build:
	./Scripts/crosscompile.sh

format:
	go fmt -x ./...

test:
	go build -o ./binaries/racetest ./Scripts/racetest.go
	./binaries/racetest --collector $(APP_ANTENNA_LISTENER_IP) --web $(API_SERVER_LISTENER_IP)
