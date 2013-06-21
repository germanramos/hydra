App Manager
===========
This application consist of 2 pieces:

### App Manager Info Server
Monitors the system resources and the port our application is listening on and notifies.

### App Manager
Ask periodically to the App Manager Info Server the state of the application and updates the corresponding Hydra server/s 

# Deploy

## Prerequisites
* build tools (gcc 4.2+, build-essential)
  * yum groupinstall 'Development Tools'
* python 2.6+
* python-devel package
  * yum install python-devel
* pip python package manager 
  * wget https://bitbucket.org/pypa/setuptools/raw/0.7.4/ez_setup.py -O - | python
  * curl -O https://raw.github.com/pypa/pip/master/contrib/get-pip.py
  * [sudo] python get-pip.py
* python psutil
  * pip install psutil 
* ssh
* git
* Increase max number of file descriptors - http://www.xenoclast.org/doc/benchmark/HTTP-benchmarking-HOWTO/node7.html


## Get Hydra source code

```
git clone https://github.com/bbva-innotech/hydra.git
```

## Configure App Manager

#### Modify ./hydra/src/app_manager/app_manager.cfg
* Modify app_id for your own app (hydra in this case).
* Set the local and cloud strategies.
* Modify the cloud name and cost.
* Add the public and private server, the public server url is the public interface for your application, the private server is where the app_manager_info_server is listening (http://127.0.0.1:7777 for example).
* Add the server_api url of one or many Hydra servers to notify (in this case could be http://localhost:7002).

# Launching

First start your application on the desired port. Then start the App Manager Info Server

```
python ./app_manager/app_manager_info_server.py <host> <public_port> <network_interface> <infoserver_port> <pid>
```

* host - the host of the application, if the App Manager Info Server is in the same machine, you can use localhost.
* public_port - port used by your application (or one of them if many are used)
* network_interface - interface the server will be started (0.0.0.0 for any interface)
* infoserver_port - port used by the App Manager Info Server to receive requests from the App Manager.
* pid - pid of your application.

Now, start the App Manager application:

```
cd ./app_manager
python ./app_manager.py
```
