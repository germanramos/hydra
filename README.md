hydra
=====

"if a head is cut off, two more will take its place"

Deploy
======
## Prerequisites
* build tools (gcc 4.2+, build-essential)
* python 2.7+
* ssh
* git

## Install nodejs
* Download nodejs source code from http://nodejs.org/download/
Latest build at the time of writing is 0.10.11

```
wget http://nodejs.org/dist/v0.10.11/node-v0.10.11.tar.gz
```
* Unzip

```
tar xvfz node-v0.10.11.tar.gz
```

* Compile and install

```
cd node-v0.10.11
./configure
make // -j<num cores + 1> for faster compiling
make install // as superuser
```

## Install mongodb
Follow this instructions: 
* Ubuntu - http://docs.mongodb.org/manual/tutorial/install-mongodb-on-ubuntu/
* CentOS/Fedora - http://docs.mongodb.org/manual/tutorial/install-mongodb-on-red-hat-centos-or-fedora-linux/

## Get Hydra source code

```
git clone https://github.com/bbva-innotech/hydra.git
```

## Install dependancies
```
cd hydra/src/lib/
npm install
```


Launching
=========

```
node ./server_api/app/main.js --port=7001 --env=pro &
node ./client_api/app/main.js --port=7002 --env=pro &
```