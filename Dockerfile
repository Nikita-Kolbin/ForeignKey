FROM golang:1.22

WORKDIR /ForeignKey

COPY ./ .

RUN go mod tidy
RUN go mod vendor

RUN go build -o build.exe ./cmd/foreign-key/main.go

CMD ["./build.exe"]
