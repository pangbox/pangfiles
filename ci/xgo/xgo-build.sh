#!/bin/sh
set -e -u -o pipefail -x

XGO_CMD=github.com/crazy-max/xgo@v0.41.0
XGO_DOCKER_IMAGE=${1:-ghcr.io/pangbox/pangfiles/xgo:latest}

mkdir -p bin /tmp/xgo-cache
go run "${XGO_CMD}" --docker-image="${XGO_DOCKER_IMAGE}" --targets=darwin/amd64 --pkg=cmd/pang --out bin/pang .
go run "${XGO_CMD}" --docker-image="${XGO_DOCKER_IMAGE}" --targets=darwin/arm64 --pkg=cmd/pang --out bin/pang .
