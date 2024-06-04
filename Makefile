run-all:
	cd deployments/docker && docker-compose up -d --force-recreate postgres kafka-ui kafka0 kafka-init-topics
	cd enrichstorage && make migrate-docker

stop-all:
	cd deployments/docker && docker-compose down