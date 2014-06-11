#!/bin/bash

### http://tecadmin.net/create-rpm-of-your-own-script-in-centosredhat/#

sudo yum install rpm-build rpmdevtools
rm -rf ~/rpmbuild
rpmdev-setuptree

mkdir ~/rpmbuild/SOURCES/hydra-3
cp ./fixtures/hydra.conf  ~/rpmbuild/SOURCES/hydra-3
cp ./fixtures/apps-example.json  ~/rpmbuild/SOURCES/hydra-3
cp hydra-init.d.sh ~/rpmbuild/SOURCES/hydra-3
cp ../../bin/hydra ~/rpmbuild/SOURCES/hydra-3

cp hydra.spec ~/rpmbuild/SPECS

pushd ~/rpmbuild/SOURCES/
tar czf hydra-3.0.tar.gz hydra-3/
cd ~/rpmbuild 
rpmbuild -ba SPECS/hydra.spec

popd
cp ~/rpmbuild/RPMS/x86_64/hydra-3-0.x86_64.rpm .
