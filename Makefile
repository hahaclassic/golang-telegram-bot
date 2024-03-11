.PHONY: build
build:
	go build -v -o ./main.exe ./cmd/main.go

build-app: 
	docker build -t link-saver:local .

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
	docker exec -it postgres psql -U admin link-saver-db

.DEFAULT_GOAL := build