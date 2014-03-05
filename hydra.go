package main

import (
	"fmt"
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
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		// fmt.Println(server.Usage() + "\n")
		// fmt.Println(err.Error() + "\n")
		log.Fatal(err.Error() + "\n")
	}

	if config.DataDir == "" {
		// TODO: include log system
		// log.Fatal("The data dir was not set and could not be guessed from machine name")
		// return errors.New("data directory attribute is required")
		fmt.Println("No data dir")
		os.Exit(1)
	}

	// Create data directory if it doesn't already exist.
	if err := os.MkdirAll(config.DataDir, 0744); err != nil {
		// log.Fatalf("Unable to create path: %s", err)
		fmt.Println("Unable to create path: %s", err)
		os.Exit(1)
	}

	// Load etcd configuration.
	if err := config.LoadEtcdConfig(); err != nil {
		fmt.Println("No load etcd config")
		os.Exit(1)
	}

	var etcd = etcd.New(config.EtcdConf)
	// etcd.Load()
	// etcd.Start()
	etcd.Load()
	go func() {
		etcd.Start()
	}()

	connector.SetEtcdConnector(etcd)
	// etcdDriver := driver.NewEtcdDriver(etcd.EtcdServer, etcd.PeerServer)
	// var server = server.NewServer(etcdDriver)
	// server.Start()

	// TODO: Use Config addr
	privateHydraListener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal("Failed to create hydra listener: ", err)
	}
	var privateServer = server.NewPrivateServer(privateHydraListener)
	privateServer.RegisterControllers()
	log.Fatal(http.Serve(privateServer.Listener, privateServer.Router))
}
