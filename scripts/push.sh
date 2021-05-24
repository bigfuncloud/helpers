#!/bin/sh
set -ex
cd "${0%/*}/.."

docker push bigfuncloud/helpers:latest
