#!/bin/bash

BRANCH=$1
TAG=$2

### Install GO 1.2+
# tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PATH:/usr/local/go/bin
cd
if [ ! -d "$HOME/go" ]; then
  mkdir go
fi
export GOPATH=$HOME/go
export GOROOT=/usr/local/go

### Install zeromq 3.2
# CentOS: http://zeromq.org/distro:centos
# Ubuntu: https://launchpad.net/~chris-lea/+archive/zeromq

### Get etcd
go get github.com/coreos/etcd
cd $GOPATH/src/github.com/coreos/etcd
git checkout v0.3.0
./build

### Get goven
go get github.com/kr/goven

### Get Hydra
go get github.com/innotech/hydra
cd /home/innotechdev/go/src/github.com/innotech/hydra
git checkout $BRANCH
cd vendors
goven github.com/coreos/etcd

### Build Hydra
# Go to hydra parent directory
cd ..
./build

### Remake tag
git push --delete origin $TAG
git tag -a $TAG
git push origin $TAG

