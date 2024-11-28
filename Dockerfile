FROM golang:1.23-bookworm

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum ./
COPY ./src ./src

RUN go mod download

CMD ["go", "run", "src/main.go"]
