FROM golang:1.15.3-alpine3.12

WORKDIR /buillding

COPY . .

RUN go build -o service ./cmd/main.go

ENTRYPOINT ["./service"]
