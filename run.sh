#!/bin/sh

set -e
go mod tidy
go vet -v ./...
go test -v ./...
