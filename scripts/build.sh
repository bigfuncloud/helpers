#!/bin/sh
set -ex
cd "${0%/*}/.."

cd biginit
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o bin/prod/biginit .

cd ..
docker build -t bigfuncloud/helpers:latest .

