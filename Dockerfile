FROM golang:1.23-alpine AS build-stage

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./src ./src

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build src/cmd/raspy-monitor/main.go

FROM alpine:3.21 AS docker-cli-builder

RUN apk add --no-cache curl tar gzip

ARG DOCKER_VERSION=28.1.1

RUN ARCH=$(uname -m) && \
    curl -SL https://download.docker.com/linux/static/stable/${ARCH}/docker-${DOCKER_VERSION}.tgz -o docker.tgz && \
    tar xzf docker.tgz --strip-components=1 docker/docker && \
    mv docker /usr/bin/docker

FROM scratch AS final-stage

COPY --from=build-stage /usr/src/app/main .
COPY --from=docker-cli-builder /usr/bin/docker /usr/bin/docker

CMD ["./main"]
