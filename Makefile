.PHONY: build
build:
	go build -o  ./bin/dickts ./app

.PHONY: run
run:
	go run ./app
