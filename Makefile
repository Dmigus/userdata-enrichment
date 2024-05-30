run-all:
	cd deployments && docker-compose up -d --force-recreate

stop-all:
	cd deployments && docker-compose down