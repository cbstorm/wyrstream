ifneq (,$(wildcard .env))
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
	$(CC) build -o $(CONTROL_SVC) $(CONTROL_SERVICE_DIR)/main.go
build-stream-svc:
	$(CC) build -o $(STREAM_SVC) $(STREAM_SERVICE_DIR)/main.go
build-auth-svc:
	$(CC) build -o $(AUTH_SVC) $(AUTH_SERVICE_DIR)/main.go $(AUTH_SERVICE_DIR)/subscribers.go
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
docker-dev:
	docker compose -f docker-compose.dev.yml up -d
mkenv:
	@echo "MONGODB_URL=\nMONGODB_DB_NAME=\nADDR=\nNATS_CORE_USERNAME=\nNATS_CORE_PASSWORD=\nNATS_CORE_HOST=\nNATS_CORE_PORT=\nNATS_CORE_QUEUE_GROUP=\n" > .env
setup: mkenv

