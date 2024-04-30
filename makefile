swagger:
	swag init -g .\cmd\foreign-key\main.go

build:
	go build -o build.exe .\cmd\foreign-key\main.go
