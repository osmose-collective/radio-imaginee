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
	-@mkdir -p logs/icecast2 data
	-@touch data/filebrowser.db data/history.txt data/latest.txt
	-@chmod 777 data logs logs/icecast2 data/filebrowser.db data/history.txt data/latest.txt
	$(ENV) docker-compose up -d --no-recreate

.PHONY: down ps build
down ps build:
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

.PHONY: re
re: down build up logs

.PHONY: cleanup-playlists
cleanup-playlists:
	@# sudo apt install faad mplayer lame ffmpeg
	find playlists -iname "*.m4a" | while read file; do (set -x; faad -o "$$file.mp3" "$$file" && rm "$$file"); sleep .1; done
	@#find playlists -iname "*.wma" | while read file; do (set -x; mplayer -vo null -vc dummy -af resample=44100 -ao pcm:waveheader "$$file" && lame -m -s audiodump.wav -o "$$file.mp3"; rm -f audiodump.wav); sleep .1; done
	@#find playlists -iname "*.wma" | while read file; do (set -x; ffmpeg -loglevel info -y -vsync 2 -i "$$file" -acodec libmp3lame -ab 128k "$$file.mp3"); sleep .1; done
