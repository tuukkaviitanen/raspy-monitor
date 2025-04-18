FROM golang:1.23-bookworm AS build-stage

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./src ./src

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build src/cmd/raspy-monitor/main.go

FROM alpine:latest AS docker-cli-builder

RUN apk add --no-cache curl tar gzip

# Specify the desired Docker CLI version
ARG DOCKER_VERSION=26.0.1

RUN curl -SL https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz -o docker.tgz && \
    tar xzf docker.tgz --strip-components=1 docker/docker && \
    mv docker /usr/local/bin/docker

FROM scratch AS final-stage

COPY --from=build-stage /usr/src/app/main .
COPY --from=docker-cli-builder /usr/local/bin/docker /usr/local/bin/docker

CMD ["./main"]
