DOCKER_COMPOSE = docker-compose -f docker/local/docker-compose.yml

.PHONY: build up down up-build

build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

app:
	$(DOCKER_COMPOSE) exec cat-backend sh
