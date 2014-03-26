#!/bin/sh -e

. ./build

# go test -i ./config
# go test -v ./config

# go test -i ./model/entity
# go test -v ./model/entity

# go test -i ./load_balancer
# go test -v ./load_balancer

# go test -i ./server
# go test -v ./server

# go test -i ./server/controller
# go test -v ./server/controller

# go test -i ./utils
# go test -v ./utils

go test -i ./tests/functional/api
HYDRA_BIN_PATH=$(pwd)/bin/hydra go test -v ./tests/functional/api

# echo "--- ETCD FUNCTIONAL TESTS ---\n"
# go test -i ./vendors/github.com/coreos/etcd/tests/functional
# ETCD_BIN_PATH=$(pwd)/bin/hydra HYDRA_ENV=ETCD_TEST go test -v ./vendors/github.com/coreos/etcd/tests/functional