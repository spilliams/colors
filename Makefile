build:
	go build -o ./bin/colors ./cmd/colors.go
.PHONY: build

test:
	go test ./...
.PHONY: test
