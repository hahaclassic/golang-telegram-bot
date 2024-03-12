.PHONY: build
build:
	go build -v -o ./main.exe ./cmd/main.go

build-app: 
	docker build -t link-saver:local .

set-env:
	export $(envsubst < ./configs/config_postgres.env)
	export $(envsubst < ./configs/config_bots.env)

run-app:
	docker run link-saver:local

build-compose:
	export PWD=$(PWD)
	docker-compose build

run-compose: 
	export PWD=$(PWD)
	docker-compose up -d

stop: 
	docker stop postgres link-saver-app

clear-postgres-data:
	sudo rm -rf ./storage/postgresql/data
	mkdir ./storage/postgresql/data
	docker rm postgres
	docker volume rm postgresql-data 

run-psql:
	docker exec -it postgres psql -U $(POSTGRES_USER) $(POSTGRES_DB)

.DEFAULT_GOAL := build