FROM golang:1.23-bookworm AS build-stage

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./src ./src

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build src/cmd/raspy-monitor/main.go

FROM scratch AS final-stage

COPY --from=build-stage /usr/src/app/main .

CMD ["./main"]
