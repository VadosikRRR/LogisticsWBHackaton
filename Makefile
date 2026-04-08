COMPOSE ?= docker compose

.PHONY: up down logs ps restart build test clean import-train import-train-dry db-patch-target2h

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f --tail=200

ps:
	$(COMPOSE) ps

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d --build

build:
	$(COMPOSE) build

test:
	cd backend && GOCACHE=/tmp/go-cache go test ./...

clean:
	$(COMPOSE) down -v --remove-orphans

import-train:
	cd backend && GOMODCACHE=/tmp/go-mod-cache GOCACHE=/tmp/go-cache go run ./importer/cmd/train_importer --file ../datasets/train_team_track.parquet

import-train-dry:
	cd backend && GOMODCACHE=/tmp/go-mod-cache GOCACHE=/tmp/go-cache go run ./importer/cmd/train_importer --file ../datasets/train_team_track.parquet --dry-run

db-patch-target2h:
	docker exec -i logistics-postgres psql -U postgres -d logistics < deploy/postgres/manual/001_add_target_2h.sql
