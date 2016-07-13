#!/usr/bin/env bash

set -e

output() {
  printf "\033[32m"
  echo $1
  printf "\033[0m"
  exit 1
}

coverage_mode=$1

test -z $coverage_mode && output "Usage: $0 coverage_mode"
test -z $(which glide) && output "glide command not found"

test -f coverage.txt && rm -rf coverage.txt
echo "mode: ${coverage_mode}" > coverage.txt

for d in $(go list ./...); do
  go test -v -race -cover -coverprofile=profile.out -covermode=${coverage_mode} $d
  if [ -f profile.out ]; then
    sed '1d' profile.out >> coverage.txt
    rm profile.out
  fi
done
