package main

import (
	// "fmt"
	"net"
	"net/http"
	"os"
	// "time"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/server"
)

func main() {
	// Load configuration.
	var conf = config.New()
	if err := conf.Load(os.Args[1:]); err != nil {
		log.Fatal(err.Error() + "\n")
	}

	if conf.DataDir == "" {
		log.Fatal("Data dir does't exist")
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(conf.DataDir, 0744); err != nil {
		log.Fatalf("Unable to create path: %s", err)
	}

	// Load etcd configuration.
	if err := conf.LoadEtcdConfig(); err != nil {
		log.Fatalf("Unable to load etcd conf: %s", err)
	}

	// Load applications.
	var appsConfig = config.NewApplicationsConfig()
	if _, err := os.Stat(conf.AppsFile); os.IsNotExist(err) {
		log.Warnf("Unable to find apps file: %s", err)
	} else {
		if err := appsConfig.Load(conf.AppsFile); err != nil {
			log.Fatalf("Unable to load applications: %s", err)
		}
	}

	var etcd = etcd.New(conf.EtcdConf)
	etcd.Load()
	hydraEnv := os.Getenv("HYDRA_ENV")
	if hydraEnv == "ETCD_TEST" {
		etcd.Start(hydraEnv)
	} else {
		go func() {
			etcd.Start(hydraEnv)
		}()

		connector.SetEtcdConnector(etcd)

		// Persist Configured applications
		if err := appsConfig.Persists(); err != nil {
			log.Fatalf("Failed to save configured applications: ", err)
		}

		// etcdDriver := driver.NewEtcdDriver(etcd.EtcdServer, etcd.PeerServer)
		// var server = server.NewServer(etcdDriver)
		// server.Start()

		// TODO: Use Config addr
		privateHydraListener, err := net.Listen("tcp", conf.PrivateAddr)
		// privateHydraListener, err := net.Listen("tcp", ":8181")
		if err != nil {
			log.Fatalf("Failed to create hydra listener: ", err)
		}
		var privateServer = server.NewPrivateServer(privateHydraListener)
		privateServer.RegisterControllers()
		log.Infof("private hydra server [name %s, listen on %s, advertised url %s]", conf.Name, conf.PrivateAddr, "http://"+conf.PrivateAddr)
		log.Fatal(http.Serve(privateServer.Listener, privateServer.Router))
	}
}
