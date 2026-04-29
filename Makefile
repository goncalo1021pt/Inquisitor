COMPOSE = docker compose

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

build:
	$(COMPOSE) build

logs:
	$(COMPOSE) logs -f

# Drop into the attacker container
attack:
	docker exec -it inquisitor /bin/bash

# Drop into the victim container
victim:
	docker exec -it inq-ftp-victim /bin/bash

# Show MAC + IP for both peers (handy when invoking inquisitor)
info:
	@echo "--- ftp-server ---"
	@docker exec inq-ftp-server sh -c "ip -4 -o addr show | awk '{print \$$2,\$$4}'; ip -o link show | awk -F'link/ether ' '/link\\/ether/ {print \$$2}' | awk '{print \$$1}'"
	@echo "--- ftp-victim ---"
	@docker exec inq-ftp-victim sh -c "ip -4 -o addr show | awk '{print \$$2,\$$4}'; ip -o link show | awk -F'link/ether ' '/link\\/ether/ {print \$$2}' | awk '{print \$$1}'"

delete:
	$(COMPOSE) down --rmi all --volumes --remove-orphans

clean: delete

re: clean up

.PHONY: up down build logs attack victim info delete clean re
