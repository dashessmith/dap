go mod tidy || exit
go generate ./... || exit
gofumpt -l -s -w . || exit
go build -o ./bin/ ./cmd/... || exit