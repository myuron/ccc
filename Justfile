default:
  @just --list

run:
  go run main.go testdata/member.csv

fmt:
  go fmt ./...

lint:
  golangci-lint run ./...
