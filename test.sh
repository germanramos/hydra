#!/bin/sh -e

. ./build

go test -i ./config
go test -v ./config

go test -i ./server
go test -v ./server