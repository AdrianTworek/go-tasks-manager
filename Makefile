build:
	go build -o bin/go-tasks-manager ./cmd

run: build
	./bin/go-tasks-manager

dev:
	air -c .air.toml

test:
	go test -v ./... -count=1