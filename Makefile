.PHONY: test
test:
	@test "$(shell echo 123)" = "$(shell echo 123)" \
		|| { echo Not equal; exit 2; } \
		&& { echo Equal; }

.PHONY: up
up: start_db run_app

.PHONY: down
down: stop_db

.PHONY: build_app
build_app:
	go build -o  ./bin/dickts ./app

.PHONY: run_app
run_app:
	go run ./app

.PHONY: start_db
start_db:
	docker-compose -f docker/Docker-compose.yaml start

.PHONY: stop_db
stop_db:
	docker-compose -f docker/Docker-compose.yaml stop

.PHONY: up_db
up_db:
	docker-compose -f docker/Docker-compose.yaml up -d

.PHONY: down_db
down_db:
	docker-compose -f docker/Docker-compose.yaml down
