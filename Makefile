build: format clean test
	go build ./...

test: get
	go test -v .

integration-test: get
	go test -v -tags integration ./...

bench: get
	go test -v -bench . ./...

get:
	go get -t -v ./...

format:
	find . -name \*.go -type f -exec gofmt -w {} \;

clean:

.PHONY: clean build
