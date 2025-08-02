NETWORK_NAME=microservices-network

SHARED_COMPOSE=docker-compose.yaml
SERVICES_DIR=services
SERVICE_COMPOSE_FILES=$(wildcard $(SERVICES_DIR)/*/docker-compose.yaml)

network:
	@docker network inspect $(NETWORK_NAME) >/dev/null 2>&1 || \
		docker network create $(NETWORK_NAME)
	@echo "Docker network '$(NETWORK_NAME)' created or already exists."

shared-up: network
	docker compose -f $(SHARED_COMPOSE) up -d
	@echo "Shared containers started."

services-up:
	@for compose_file in $(SERVICE_COMPOSE_FILES); do \
		echo "Starting service: $$compose_file"; \
		docker compose -f $$compose_file up -d; \
	done

up: shared-up services-up
	@echo "All containers are up and running."

down:
	-docker compose -f $(SHARED_COMPOSE) down
	@for compose_file in $(SERVICE_COMPOSE_FILES); do \
		echo "Stopping service: $$compose_file"; \
		docker compose -f $$compose_file down; \
	done
	@echo "All containers stopped."

clean: down
	@docker image prune -f
	@echo "Docker images cleaned."

rebuild: clean up
	@echo "Full rebuild complete."

.PHONY: network shared-up services-up up down clean rebuild
