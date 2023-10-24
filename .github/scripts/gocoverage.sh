#!/bin/sh

unit_coverage_test() {
  go mod download
  go test -race -covermode atomic -coverprofile=coverage.out ./...
}

unit_coverage_test