test-ci:
	go test ./...

test-coverage:
	go tool cover -func=coverage.out