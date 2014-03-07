package main

import (
	"net"
	"net/http"
	"os"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/server"
)

func main() {
	// Load configuration.
	var config = config.New()
	if err := config.Load(os.Args[1:]); err != nil {
		log.Fatal(err.Error() + "\n")
	}

	if config.DataDir == "" {
		log.Fatal("Data dir does't exist")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(config.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Load etcd configuration.
	if err := config.LoadEtcdConfig(); err != nil {
		log.Fatalf("Unable to load etcd config: %s", err)
	}

	var etcd = etcd.New(config.EtcdConf)
	etcd.Load()
	hydraEnv := os.Getenv("HYDRA_ENV")
	go func() {
		etcd.Start(hydraEnv)
	}()

	connector.SetEtcdConnector(etcd)
	// etcdDriver := driver.NewEtcdDriver(etcd.EtcdServer, etcd.PeerServer)
	// var server = server.NewServer(etcdDriver)
	// server.Start()

	// TODO: Use Config addr
	// privateHydraListener, err := net.Listen("tcp", config.PrivateAddr)
	privateHydraListener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to create hydra listener: ", err)
	}
	var privateServer = server.NewPrivateServer(privateHydraListener)
	privateServer.RegisterControllers()
	log.Infof("private hydra server [name %s, listen on %s, advertised url %s]", config.Name, config.PrivateAddr, "http://"+config.PrivateAddr)
	log.Fatal(http.Serve(privateServer.Listener, privateServer.Router))
}
