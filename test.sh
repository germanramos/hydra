#!/bin/sh -e

. ./build

# go test -i ./config
# go test -v ./config

# go test -i ./server
# go test -v ./server

go test -i ./model/entity
go test -v ./model/entity

# go test -i ./utils
# go test -v ./utils

# go test -i ./tests
# HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests

# go test -i ./tests/functional
# HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests/functional

# go test -i ./tests/functional/api/application
# HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests/functional/api/application