.PHONY: fresh

SHELL := /bin/bash

fresh:
	ENV=LOCAL \
	PORT=3000 \
	VERSION=VERSION \
	DB_DSN="host=localhost user=postgres password=postgres dbname=BOARD port=5432 sslmode=disable TimeZone=Asia/Shanghai" \
	fresh -c ./cmd/runner.conf
