.PHONY: run-all
run-all:
	cd deployments/docker && docker-compose up -d postgres kafka-ui kafka0 kafka-init-topics
	cd enrichstorage && make migrate-docker
	cd enricher && go build -o ./bin/app ./cmd
	cd enrichstorage && go build -o ./bin/app ./cmd
	cd deployments/docker && docker-compose build -q
	cd deployments/docker && docker-compose up -d --force-recreate enricher enrichstorage

.PHONY: stop-all
stop-all:
	cd deployments/docker && docker-compose down