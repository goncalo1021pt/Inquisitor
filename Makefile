up:
	docker compose  up -d

down:
	docker compose  down

build:
	docker compose  build

delete:
	docker compose  down --rmi all --volumes --remove-orphans

clean: delete

re: clean up

.PHONY: up down build delete clean re