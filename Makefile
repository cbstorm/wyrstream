ifneq (,$(wildcard ./.env))
    include .env
    export
endif
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
	go build -o $(CONTROL_SVC) $(CONTROL_SERVICE_DIR)
build-stream-svc:
	go build -o $(STREAM_SVC) $(STREAM_SERVICE_DIR)
build-auth-svc:
	go build -o $(AUTH_SVC) $(AUTH_SERVICE_DIR)
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

work:
	go work init $(CONTROL_SERVICE_DIR) $(STREAM_SERVICE_DIR) $(AUTH_SERVICE_DIR) $(LIB_DIR)
up:
	docker compose -f docker-compose.dev.yml up -d
down:
	docker compose -f docker-compose.dev.yml down
clean-docker:
	rm -rf .opt
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
test:
	cd lib && go test ./... && cd ..
	cd auth_service && go test ./... && cd ..
	cd control_service && go test ./... && cd ..
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
publish:
	ffmpeg \
	-v error \
	-re \
	-stream_loop -1 \
	-i tmp/vid_0.mp4 \
    -maxrate:v 4096k \
    -bufsize:v 1024k \
    -preset ultrafast \
	-f mpegts "srt://127.0.0.1:6000?streamid=publish:/live/STR66E95B8E2?key=vf5ISSbAo20E4pjgJnuAHWQvggtGtF"

