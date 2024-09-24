ifneq (,$(wildcard ./.env))
    include .env
    export
endif
CONTROL_SERVICE_DIR = ./control_service
STREAM_SERVICE_DIR = ./stream_service
AUTH_SERVICE_DIR = ./auth_service
HLS_SERVICE_DIR = ./hls_service
ALERT_SERVICE_DIR = ./alert_service
LIB_DIR = ./lib
OUT = ./dist
CONTROL_SVC = $(OUT)/control_svc
STREAM_SVC = $(OUT)/stream_svc
AUTH_SVC = $(OUT)/auth_svc
HLS_SVC = $(OUT)/hls_svc

all: control-svc stream-svc auth-svc

build: mkdist build-control-svc build-stream-svc build-auth-svc build-hls-svc build-alert-svc
build-control-svc:
	go build -o $(CONTROL_SVC) $(CONTROL_SERVICE_DIR)
build-stream-svc:
	go build -o $(STREAM_SVC) $(STREAM_SERVICE_DIR)
build-auth-svc:
	go build -o $(AUTH_SVC) $(AUTH_SERVICE_DIR)
build-hls-svc:
	go build -o $(AUTH_SVC) $(HLS_SERVICE_DIR)
build-alert-svc:
	cd $(ALERT_SERVICE_DIR) && yarn build && cd ..
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
	rm -rf tmp/*.jpg

work:
	rm -rf go.work
	go work init $(CONTROL_SERVICE_DIR) $(STREAM_SERVICE_DIR) $(AUTH_SERVICE_DIR) $(LIB_DIR) $(HLS_SERVICE_DIR)
	go work sync
deps:
	cd $(LIB_DIR) && go mod download && cd ..
	cd $(AUTH_SERVICE_DIR) && go mod download && cd ..
	cd $(CONTROL_SERVICE_DIR) && go mod download && cd ..
	cd $(STREAM_SERVICE_DIR) && go mod download && cd ..
	cd $(ALERT_SERVICE_DIR) && yarn install && cd ..
up:
	docker compose -f docker-compose.dev.yml up -d
down:
	docker compose -f docker-compose.dev.yml down
clean-docker:
	rm -rf .opt
mkenv:
	@echo "MONGODB_URL=\n \
	MONGODB_DB_NAME=\n \
	STREAM_SERVER_ADDR=\n \
	STREAM_SERVER_PUBLIC_URL=\n \
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
	DISCORD_BOT_TOKEN=\n \
	MINIO_HOST=\n \
	MINIO_PORT=\n \
	MINIO_ACCESS_KEY=\n \
	MINIO_SECRET_KEY=\n \
	MINIO_BUCKET_NAME=\n \
	MINIO_PUBLIC_URL=\n \
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
	-f mpegts "srt://127.0.0.1:6000?streamid=publish:/live/STR66F24EE21?key=JYqXgwdxlUkD2yXqN6L2TGTygYzpkN"
test-sub:
	ffplay -v quiet -f mpegts -transtype live -i "srt://127.0.0.1:6000?streamid=/live/STR66E95B8E2?key=0MRWUlRLHSViEddcmOtKLMDYann1st"
test-play-hls:
	ffplay -i "http://127.0.0.1:9000/wyrstream/streams/STR66F24EE21/playlist/playlist.m3u8"
test-hls:
	ffmpeg \
	-i srt://127.0.0.1:6000?streamid=/live/STR66E95B8E2?key=0MRWUlRLHSViEddcmOtKLMDYann1st \
	-c:v libx264 \
	-c:a aac \
	-b:a 160k \
	-b:v 1M \
	-maxrate:v 1M \
	-bufsize 512k \
	-crf 18 \
	-preset ultrafast \
	-f hls \
	-hls_time 6 \
	-hls_list_size 6 \
	-hls_segment_filename hls_service/public/STR66E95B8E2/seg-%05d.ts \
	-start_number 1 \
	hls_service/public/STR66E95B8E2/playlist.m3u8
test-thumbnail:
	ffmpeg -v error -i hls_service/public/STR66EC7A661/seg-00004.ts -q:v 1 -frames:v 1 -y tmp/STR66EC7A942.jpg