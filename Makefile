run:
	go run ./cmd

build:
	go build -o ./bin/reboot ./cmd/reboot/main.go

build-c:
	CC=musl-gcc go build -ldflags '-linkmode external -extldflags "-static"' -o ./bin/reboot ./cmd/reboot/main.go

docker-build:
	docker build -f Dockerfile.reboot -t inneroot/reboot:latest --platform linux/amd64 .
