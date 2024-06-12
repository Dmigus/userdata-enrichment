
run-storage:
	cd deployments/docker && docker-compose up --wait postgres
	cd enrichstorage && make migrate-docker

run-kafka:
	cd deployments/docker && docker-compose up -d kafka-ui kafka0 kafka-init-topics

.PHONY: run-all
run-all: run-storage run-kafka
	cd enricher && go build -o ./bin/app ./cmd
	cd enrichstorage && go build -o ./bin/app ./cmd
	cd deployments/docker && docker-compose build -q
	cd deployments/docker && docker-compose up -d --force-recreate enricher enrichstorage

.PHONY: stop-all
stop-all:
	cd deployments/docker && docker-compose down