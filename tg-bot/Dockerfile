FROM golang:1.22

WORKDIR /tg-bot

COPY . .

RUN go mod tidy
RUN go mod vendor

RUN go build -o build ./cmd/main.go

CMD ["./build"]
