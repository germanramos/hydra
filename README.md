Hydra
=====

"If a head is cut off, two more will take its place"

# What is Hydra?
Hydra is multi-cloud application discovery, management and balancing service.

Hydra attempts to ease the routing and balancing burden from servers and delegate it on the client (browser, mobile app, etc).   

Deploy
======
## Prerequisites
* build tools (gcc 4.2+, build-essential)
* python 2.7+
* python-devel package
* pip python package manager
* python psutil (install via pip)
* ssh
* git
* Increase max number of file descriptors - http://www.xenoclast.org/doc/benchmark/HTTP-benchmarking-HOWTO/node7.html

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

## Configure Hydra

### Configure Hydra Server

#### Modify ./hydra/src/lib/config/local.json
* Change url for the public url of your service, this is the url that will be used for other hydras.
* Change server port, this is the base port that will be used to sync with other hydras.
* Modify timeouts on app and server.
* Modify mongodb configuration.
* Add your own QLog credentials.

#### Modify ./hydra/src/app_manager/app_manager.cfg
* Modify app_id for your own app (hydra in this case).
* Set the local and cloud strategies.
* Modify the cloud name and cost.
* Add the public and private server, the public server url is the client_api (same as in local.json), the private server is where the app_manager_info_server is listening (http://127.0.0.1:7777 for example).
* Add the server_api url of an Hydra to notify (in this case could be http://localhost:7002).

Launching
=========

Use the command run.sh on the src folder

```
./run.sh <envirnoment> <client_api_port> <server_api_port> <app_manager_info_server_port>
```

* environment - name of the environment file to use (local)
* client_api_port - port used by the client api
* server_api_port - port used by the server api and syncronization
* app_manager_info_server - port used by the app_manager_info_server to monitor the system

Alternatively, there is an startup script that can be copied to /etc/init.d/ as superuser once configured the desired parameters:

```
cp hydra/src/hydra /etc/init.d/
```

or used as is:
```
./hydra start
./hydra stop
```

Using main ports, such as 80 or 443, requires that no other application is listening on that port and superuser rights.
