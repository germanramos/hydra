#!/bin/sh -e

. ./build

go test -i ./config
go test -v ./config