#!/bin/bash

### http://linuxconfig.org/easy-way-to-create-a-debian-package-and-local-package-repository

rm -rf ~/debbuild
mkdir -p ~/debbuild/DEBIAN
cp control ~/debbuild/DEBIAN

mkdir -p ~/debbuild/etc/hydra
cp ../hydra.conf ~/debbuild/etc/hydra
cp ../apps-example.json ~/debbuild/etc/hydra

mkdir -p ~/debbuild/etc/init.d
cp hydra-basic-probe-init.d.sh ~/debbuild/etc/init.d/hydra

mkdir -p ~/debbuild/usr/local/hydra
cp ../hydra  ~/debbuild/usr/local/hydra
# cp ../src/hydra-basic-probe.py  ~/debbuild/usr/local/hydra
# cp ../src/parseStatusDat.py  ~/debbuild/usr/local/hydra

chmod -R 644 ~/debbuild/usr/local/hydra/* ~/debbuild/etc/hydra/*
chmod 755 ~/debbuild/etc/init.d/hydra
chmod 755 ~/debbuild/usr/local/hydra/hydra

adduser --no-create-home --disabled-login hydra
sudo chown -R hydra:hydra ~/debbuild/*

pushd ~
sudo dpkg-deb --build debbuild

popd
sudo mv ~/debbuild.deb hydra.noarch.deb
