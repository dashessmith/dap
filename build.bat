go mod tidy || exit
@rem gofmt -s -w .
gofumpt -l -s -w . || exit
go build -o  . ./cmd/... || exit