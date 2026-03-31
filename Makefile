.PHONY: build test clean

build:
	go build -o extractor ./cmd/extractor

test:
	go test ./...

clean:
	rm -f extractor

cross-build:
	GOOS=darwin GOARCH=amd64 go build -o extractor-darwin-amd64 ./cmd/extractor
	GOOS=darwin GOARCH=arm64 go build -o extractor-darwin-arm64 ./cmd/extractor
	GOOS=linux GOARCH=amd64 go build -o extractor-linux-amd64 ./cmd/extractor
	GOOS=windows GOARCH=amd64 go build -o extractor-windows-amd64.exe ./cmd/extractor