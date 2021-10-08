include chicha.conf
.DEFAULT_GOAL := start

start:
	go run ./chicha.go

build:
	./Scripts/crosscompile.sh

test:
	./Scripts/RaceTest.sh

