# Hydra Installation

## Ubuntu/Debian

Add PPAs for:  
https://launchpad.net/~chris-lea/+archive/libpgm  
https://launchpad.net/~chris-lea/+archive/zeromq  
  
and execute:  
```
sudo dpkg -i hydra-3-0.x86_64.deb
sudo apt-get install -f
```
## CentOS/RedHat/Fedora
```
sudo yum install libzmq3-3.2.2-13.1.x86_64.rpm hydra-3-0.x86_64.rpm
```

## Run
Check the configurations files "hydra.conf" and "apps.json" under /etc/hydra directory. After that just run:
```
sudo /etc/init.d/hydra start
```

# Hydra Configuration

The hydra configuration file ("/etc/hydra/hydra.conf") includes many of the options availables in the etcd configuration and these arguments preserve the same names than the original version of etcd.

Configuration options can be set in two places:

 1. Command line flags
 2. Configuration file

Options set on the command line take precedence over all other sources.

## Command Line Flags

### Required

* `-name` - The node name. Defaults to the hostname.

### Optional

#### Etcd arguments
* `-addr` - The advertised public hostname:port for client communication. Defaults to `127.0.0.1:7401`.
* `-apps-file` - The path of the application configuration file. Defaults to `/etc/hydra/apps.json`.
* `-bind-addr` - The listening hostname for client communication. Defaults to advertised IP.
* `-discovery` - A URL to use for discovering the peer list. (i.e `"https://discovery.etcd.io/your-unique-key"`).
* `-peers` - A comma separated list of peers in the cluster (i.e `"203.0.113.101:7701,203.0.113.102:7701"`).
* `-ca-file` - The path of the client CAFile. Enables client cert authentication when present.
* `-cert-file` - The cert file of the client.
* `-key-file` - The key file of the client.
* `-config` - The path of the etcd configuration file. Defaults to `/etc/hydra/hydra.conf`.
* `-data-dir` - The directory to store log and snapshot. Defaults to the current working directory.
* `-f, -force` - The node is started as a standalone server when it can not join the cluster.
* `-peer-addr` - The advertised public hostname:port for server communication. Defaults to `127.0.0.1:7701`.
* `-peer-bind-addr` - The listening hostname for server communication. Defaults to advertised IP.
* `-peer-ca-file` - The path of the CAFile. Enables client/peer cert authentication when present.
* `-peer-cert-file` - The cert file of the server.
* `-peer-heartbeat-timeout` - This is the frequency with which the leader will notify followers that it is still the leader and this is also a delay for how long it takes for commands to be committed. Default to 50 milliseconds.
* `-peer-election-timeout` - This timeout is how long a follower node will go without hearing a heartbeat before attempting to become leader itself. Default to 200 milliseconds.
* `-peer-key-file` - The key file of the server.
* `-snapshot=false` - Disable log snapshots. Defaults to `true`.
* `-snapshot-count` - Time interval in milliseconds between the log snapshot are made.

#### Hydra arguments
* `-instance-expiration-time` - This is the ttl for instance information.
* `-private-addr` - The hydra private api hostname:port for probe communication. Defaults to `127.0.0.1:7771`.
* `-public-addr` - The hydra public api hostname:port for client communication. Defaults to `127.0.0.1:7772`.
* `-load-balancer-addr` - The hydra load balancer hostname:port for internal and worker communication. Defaults to `*:7777`.
* `-v, -verbose` - Show logs in DEBUG mode. Defaults to `false`

## Configuration File

The hydra configuration file is written in [TOML](https://github.com/mojombo/toml)
and read from `/etc/hydra/hydra.conf` by default.

```TOML
addr = "127.0.0.1:7401"
apps_file = ""
bind_addr = ""
ca_file = ""
cert_file = ""
data_dir = "."
discovery = "http://etcd.local:4001/v2/keys/_etcd/registry/examplecluster"
force = false
instance_expiration_time = 300
key_file = ""
load_balancer_addr = "*:7777"
peers = []
private_addr = "127.0.0.1:7771"
public_addr = "127.0.0.1:7772"
name = "default-name"
snapshot = true
snapshot_count = 2000
verbose = false

[peer]
addr = "127.0.0.1:7701"
bind_addr = ""
ca_file = ""
cert_file = ""
key_file = ""
heartbeat_timeout = 100
election_timeout = 400
```

# Applications Configuration

The applications settings ("/etc/init.d/apps.json") allow to define for each application how they will be balanced by specifying the (balancers) workers that will be part of the chain and their arguments. For example:


```JSON
[{
	"App1": {
		"Balancers": [
			{
				"worker": "MapAndSort",
				"mapAttr": "cloud",
				"mapSort": ["google", "amazon", "azure"]
			},
			{
				"worker": "SortByNumber",
				"sortAttr": "cpuLoad",
				"order": 1
			}
		}
	}
}, {
	"App2": {
		"Balancers": {
			{
				"worker": "MapByLimit",
				"limitAttr": "limit",
				"limitValue": 50,
				"mapSort": "reverse"
			},
			{
				"worker": "RoundRobin"
				"simple": "OK"
			}
		}
	}
}]
```
