-include config.mk

ADMIN_PASSWORD	?=	secure
HARBOR_PASSWORD	?=	secure
SOURCE_PASSWORD	?=	secure
RELAY_PASSWORD  ?=	secure
MYSQL_PASSWORD  ?=	secure
PIWIK_PASSWORD  ?=	secure

ENV ?=			HARBOR_PASSWORD=$(HARBOR_PASSWORD) \
			LIVE_PASSWORD=$(HARBOR_PASSWORD) \
			ICECAST_SOURCE_PASSWORD=$(SOURCE_PASSWORD) \
			ICECAST_ADMIN_PASSWORD=$(ADMIN_PASSWORD) \
			ICECAST_PASSWORD=$(ADMIN_PASSWORD) \
			ICECAST_RELAY_PASSWORD=$(RELAY_PASSWORD) \

.PHONY: all
all: up logs

.PHONY: re-main
re-main: up
	$(ENV) docker-compose up -d --no-deps --force-recreate main

.PHONY: re-controller
re-controller: up
	$(ENV) docker-compose up -d --build --no-deps --force-recreate controller

.PHONY: up
up:
	-@mkdir -p logs/icecast2 data; chmod 777 data logs logs/icecast2; true
	$(ENV) docker-compose up -d --no-recreate

.PHONY: down ps
down ps:
	$(ENV) docker-compose $@

.PHONY: logs
logs:
	$(ENV) docker-compose logs --tail=100 -f

.PHONY: telnet
telnet:
	nc -v localhost 5000

.PHONY: skip
skip:
	echo "main(dot)harbor.skip\r\nexit" | nc -v localhost 5000
