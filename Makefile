
.PHONY: run-storage-docker
run-storage-docker:
	cd deployments/docker && docker-compose up --wait  --force-recreate postgres
	cd enrichstorage && make migrate-docker

.PHONY: run-kafka-docker
run-kafka-docker:
	cd deployments/docker && docker-compose up -d --force-recreate kafka-ui kafka0 kafka-init-topics

.PHONY: run-rabbit-docker
run-rabbit-docker:
	cd deployments/docker && docker-compose up -d --force-recreate rabbitmq

.PHONY: run-keycloak-docker
run-keycloak-docker:
	cd deployments/docker && docker-compose up -d --wait  --force-recreate  keycloak keycloak-init

.PHONY: run-all-docker
run-all-docker: run-storage-docker run-rabbit-docker run-keycloak-docker
	cd enricher && go build -o ./bin/app ./cmd
	cd enrichstorage && go build -o ./bin/app ./cmd
	cd deployments/docker && docker-compose build -q
	cd deployments/docker && docker-compose up -d --force-recreate enricher enrichstorage

.PHONY: stop-all-docker
stop-all-docker:
	cd deployments/docker && docker-compose down