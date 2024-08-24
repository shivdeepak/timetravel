.PHONY: build run test clean

build: fmt tidy
	go build -o bin/timetravel

tidy:
	go mod tidy

run: build
	./bin/timetravel

test:
	go test ./...

clean:
	go clean
	rm -rf bin/

fmt:
	go fmt ./...
