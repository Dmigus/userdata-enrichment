
.PHONY: run-storage-docker
run-storage-docker:
	cd deployments/docker && docker-compose up --wait postgres
	cd enrichstorage && make migrate-docker

.PHONY: run-kafka-docker
run-kafka-docker:
	cd deployments/docker && docker-compose up -d kafka-ui kafka0 kafka-init-topics

.PHONY: run-all-docker
run-all-docker: run-storage-docker run-kafka-docker
	cd enricher && go build -o ./bin/app ./cmd
	cd enrichstorage && go build -o ./bin/app ./cmd
	cd deployments/docker && docker-compose build -q
	cd deployments/docker && docker-compose up -d --force-recreate enricher enrichstorage

.PHONY: stop-all-docker
stop-all-docker:
	cd deployments/docker && docker-compose down