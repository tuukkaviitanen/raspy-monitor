FROM golang:1.23-bookworm

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./src ./src

RUN go build src/cmd/raspy-monitor/main.go

CMD ["./main"]
