#!/bin/bash

### http://tecadmin.net/create-rpm-of-your-own-script-in-centosredhat/#

sudo yum install rpm-build rpmdevtools
rm -rf ~/rpmbuild
rpmdev-setuptree

mkdir ~/rpmbuild/SOURCES/hydra
cp ../hydra.conf  ~/rpmbuild/SOURCES/hydra
cp ../apps-example.json  ~/rpmbuild/SOURCES/hydra
cp hydra-init.d.sh ~/rpmbuild/SOURCES/hydra

cp hydra.spec ~/rpmbuild/SPECS

pushd ~/rpmbuild/SOURCES/
tar czf hydra.tar.gz hydra/
cd ~/rpmbuild 
rpmbuild -ba SPECS/hydra.spec

popd
cp ~/rpmbuild/RPMS/noarch/hydra.noarch.rpm .
