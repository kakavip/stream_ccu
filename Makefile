#!make
export $(shell sed 's/=.*//' docker-compose.local.env)
args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

build:
	COMPOSE_PROJECT_NAME=vimai_local docker-compose -f docker-compose.local.yml build $(call args," ")

up:
	COMPOSE_PROJECT_NAME=vimai_local docker-compose -f docker-compose.local.yml up -d $(call args," ")

down:
	COMPOSE_PROJECT_NAME=vimai_local docker-compose -f docker-compose.local.yml down