ifneq (,$(wildcard ./.env))
    include .env
    export
endif
CC = go
CONTROL_SERVICE_DIR = ./control_service
STREAM_SERVICE_DIR = ./stream_service
AUTH_SERVICE_DIR = ./auth_service
LIB_DIR = ./lib
OUT = ./dist
CONTROL_SVC = $(OUT)/control_svc
STREAM_SVC = $(OUT)/stream_svc
AUTH_SVC = $(OUT)/auth_svc

all: control-svc stream-svc auth-svc

build: mkdist build-control-svc build-stream-svc build-auth-svc
build-control-svc:
	$(CC) build -o $(CONTROL_SVC) $(CONTROL_SERVICE_DIR)
build-stream-svc:
	$(CC) build -o $(STREAM_SVC) $(STREAM_SERVICE_DIR)
build-auth-svc:
	$(CC) build -o $(AUTH_SVC) $(AUTH_SERVICE_DIR)
control-svc:
	$(CONTROL_SVC)
stream-svc:
	$(STREAM_SVC)
auth-svc:
	$(AUTH_SVC)

mkdist:
	mkdir -p $(OUT)

clean:
	rm -rf $(OUT)/*
	rm -rf go.work

work:
	$(CC) work init $(CONTROL_SERVICE_DIR) $(STREAM_SERVICE_DIR) $(AUTH_SERVICE_DIR) $(LIB_DIR)
up:
	docker compose -f docker-compose.dev.yml up -d
down:
	docker compose -f docker-compose.dev.yml down
mkenv:
	@echo "MONGODB_URL=\n \
	MONGODB_DB_NAME=\n \
	ADDR=\n \
	NATS_CORE_USERNAME=\n \
	NATS_CORE_PASSWORD=\n \
	NATS_CORE_HOST=\n \
	NATS_CORE_PORT=\n \
	NATS_CORE_QUEUE_GROUP=\n \
	HTTP_HOST=\n \
	HTTP_PORT=\n \
	REDIS_USERNAME=\n \
	REDIS_PASSWORD=\n \
	REDIS_HOST=\n \
	REDIS_PORT=\n \
	REDIS_KEY_PREFIX=\n \
	> .env
setup: mkenv
route:
	npx plop route
entity:
	npx plop entity
service:
	npx plop service
repository:
	npx plop repository
dto:
	npx plop dto
module:
	npx plop module

