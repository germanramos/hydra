#!/bin/sh -e

. ./build

# go test -i ./config
# go test -v ./config

# go test -i ./server
# go test -v ./server

# go test -i ./tests
# HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests
go test -i ./tests/functional
HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests/functional