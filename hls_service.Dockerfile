FROM golang:latest AS build
WORKDIR /wyrstream
COPY hls_service ./hls_service
COPY lib ./lib
RUN go work init ./hls_service ./lib
RUN go work sync
RUN go mod download
RUN go build -v -o /wyrstream/bin/hls_service ./hls_service

FROM linuxserver/ffmpeg:latest
WORKDIR /
RUN apt-get update
RUN apt-get install -y ca-certificates

COPY --from=build /wyrstream/bin /wyrstream/bin

RUN ffmpeg -version
ENTRYPOINT [ "/wyrstream/bin/hls_service" ]
