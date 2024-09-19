ifneq (,$(wildcard ./.env))
    include .env
    export
endif
CONTROL_SERVICE_DIR = ./control_service
STREAM_SERVICE_DIR = ./stream_service
AUTH_SERVICE_DIR = ./auth_service
HLS_SERVICE_DIR = ./hls_service
LIB_DIR = ./lib
OUT = ./dist
CONTROL_SVC = $(OUT)/control_svc
STREAM_SVC = $(OUT)/stream_svc
AUTH_SVC = $(OUT)/auth_svc
HLS_SVC = $(OUT)/hls_svc

all: control-svc stream-svc auth-svc

build: mkdist build-control-svc build-stream-svc build-auth-svc build-hls-svc
build-control-svc:
	go build -o $(CONTROL_SVC) $(CONTROL_SERVICE_DIR)
build-stream-svc:
	go build -o $(STREAM_SVC) $(STREAM_SERVICE_DIR)
build-auth-svc:
	go build -o $(AUTH_SVC) $(AUTH_SERVICE_DIR)
build-hls-svc:
	go build -o $(AUTH_SVC) $(HLS_SERVICE_DIR)
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
	rm -rf $(HLS_SERVICE_DIR)/public/*

work:
	rm -rf go.work
	go work init $(CONTROL_SERVICE_DIR) $(STREAM_SERVICE_DIR) $(AUTH_SERVICE_DIR) $(LIB_DIR) $(HLS_SERVICE_DIR)
	go work sync
deps:
	cd $(LIB_DIR) && go mod tidy && cd ..
	cd $(AUTH_SERVICE_DIR) && go mod tidy && cd ..
	cd $(CONTROL_SERVICE_DIR) && go mod tidy && cd ..
	cd $(STREAM_SERVICE_DIR) && go mod tidy && cd ..
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
	PUBLIC_URL=\n \
	NATS_CORE_USERNAME=\n \
	NATS_CORE_PASSWORD=\n \
	NATS_CORE_HOST=\n \
	NATS_CORE_PORT=\n \
	NATS_CORE_QUEUE_GROUP=\n \
	HTTP_HOST=\n \
	HTTP_PORT=\n \
	HLS_HTTP_HOST=\n \
	HLS_HTTP_PORT=\n \
	HLS_PUBLIC_URL=\n \
	REDIS_USERNAME=\n \
	REDIS_PASSWORD=\n \
	REDIS_HOST=\n \
	REDIS_PORT=\n \
	REDIS_KEY_PREFIX=\n \
	> .env
setup: mkenv
test:
	cd $(LIB_DIR) && go test ./... && cd ..
	cd $(AUTH_SERVICE_DIR) && go test ./... && cd ..
	cd $(CONTROL_SERVICE_DIR) && go test ./... && cd ..
	cd $(STREAM_SERVICE_DIR) && go test ./... && cd ..
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
test-pub:
	ffmpeg \
	-v error \
	-re \
	-stream_loop -1 \
	-i tmp/vid_0.mp4 \
	-c:v libx264 \
	-b:v 2M \
    -maxrate:v 2M \
    -bufsize:v 1M \
    -preset ultrafast \
	-f mpegts "srt://127.0.0.1:6000?streamid=publish:/live/STR66E95B8E2?key=vf5ISSbAo20E4pjgJnuAHWQvggtGtF"
test-sub:
	ffplay -v quiet -f mpegts -transtype live -i "srt://127.0.0.1:6000?streamid=/live/STR66E95B8E2?key=0MRWUlRLHSViEddcmOtKLMDYann1st"
test-play-hls:
	ffplay -i "http://127.0.0.1:10000/STR66E95B8E2/playlist.m3u8"
test-hls:
	ffmpeg \
	-i srt://127.0.0.1:6000?streamid=/live/STR66E95B8E2?key=0MRWUlRLHSViEddcmOtKLMDYann1st \
	-c:v libx264 \
	-c:a aac \
	-b:a 160k \
	-b:v 2M \
	-maxrate:v 2M \
	-bufsize 1M \
	-crf 18 \
	-preset ultrafast \
	-f hls \
	-hls_time 6 \
	-hls_list_size 6 \
	-hls_segment_filename hls_service/public/STR66E95B8E2/seg-%05d.ts \
	-start_number 1 \
	hls_service/public/STR66E95B8E2/playlist.m3u8

