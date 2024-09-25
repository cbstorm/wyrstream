FROM golang:latest AS build
WORKDIR /wyrstream
COPY auth_service ./auth_service
COPY lib ./lib
RUN go work init ./auth_service ./lib
RUN go work sync
RUN go mod download
RUN go build -v -o /wyrstream/bin/auth_service ./auth_service

FROM ubuntu:latest

WORKDIR /
RUN apt-get update
RUN apt-get install -y ca-certificates

COPY --from=build /wyrstream/bin /wyrstream/bin

CMD ["/wyrstream/bin/auth_service"]
