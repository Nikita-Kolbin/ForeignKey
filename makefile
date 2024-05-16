swagger:
	swag init -g .\cmd\foreign-key\main.go

build:
	go build -o build.exe .\cmd\foreign-key\main.go

run: swagger
	go run .\cmd\foreign-key\main.go

docker_build: swagger
	docker build . -t foreign_key

docker_run: docker_build
	docker run -d -p 8082:8082 foreign_key
