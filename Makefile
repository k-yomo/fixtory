
.PHONY: test
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: generate
generate:
	go run cmd/fixtory/main.go -type=Author,Article -output=example/article.fixtory.go example