#!/bin/bash

### http://linuxconfig.org/easy-way-to-create-a-debian-package-and-local-package-repository

rm -rf ~/debbuild
mkdir -p ~/debbuild/DEBIAN
cp control ~/debbuild/DEBIAN

mkdir -p ~/debbuild/etc/hydra
cp ./fixtures/hydra.conf ~/debbuild/etc/hydra
cp ./fixtures/apps-example.json ~/debbuild/etc/hydra

mkdir -p ~/debbuild/etc/init.d
cp hydra-init.d.sh ~/debbuild/etc/init.d/hydra-core

mkdir -p ~/debbuild/usr/local/hydra
cp ../../bin/hydra  ~/debbuild/usr/local/hydra

chmod -R 644 ~/debbuild/usr/local/hydra/* ~/debbuild/etc/hydra/*
chmod 755 ~/debbuild/etc/init.d/hydra-core
chmod 755 ~/debbuild/usr/local/hydra/hydra

sudo chown -R root:root ~/debbuild/*

pushd ~
sudo dpkg-deb --build debbuild

popd
sudo mv ~/debbuild.deb hydra-core-3-0.x86_64.deb
