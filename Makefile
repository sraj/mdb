PWD := $(shell pwd)

build:
	@docker-compose build

start:
	@docker-compose up -d

stop:
	@docker-compose stop

status:
	@docker-compose ps

test:
	@docker-compose run --rm mdb go test -v

mongocli:
	@docker-compose exec mongo sh

tail:
	@docker-compose logs -f

clean: stop
	@docker-compose rm --force

.PHONY: build start stop status test mongocli tail clean