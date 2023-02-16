include chicha.conf
.DEFAULT_GOAL := start

start:
	go run ./chicha.go

build:
	./Scripts/crosscompile.sh

format:
	go fmt -x ./...

test:
	go run ./Scripts/racetest.go

