go mod tidy
gofmt -s -w .
go build -o ./bin/ ./cmd/...