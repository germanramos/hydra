#!/bin/bash

### http://tecadmin.net/create-rpm-of-your-own-script-in-centosredhat/#

sudo yum install rpm-build rpmdevtools
rm -rf ~/rpmbuild
rpmdev-setuptree

mkdir ~/rpmbuild/SOURCES/hydra-core-3
cp ./fixtures/hydra.conf  ~/rpmbuild/SOURCES/hydra-core-3
cp ./fixtures/apps-example.json  ~/rpmbuild/SOURCES/hydra-core-3
cp hydra-init.d.sh ~/rpmbuild/SOURCES/hydra-core-3
cp ../../bin/hydra ~/rpmbuild/SOURCES/hydra-core-3

cp hydra-core.spec ~/rpmbuild/SPECS

pushd ~/rpmbuild/SOURCES/
tar czf hydra-core-3.0.tar.gz hydra-core-3/
cd ~/rpmbuild 
rpmbuild -ba SPECS/hydra-core.spec

popd
cp ~/rpmbuild/RPMS/x86_64/hydra-core-3-0.x86_64.rpm .
