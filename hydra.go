package main

import (
	"fmt"
	"os"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/driver"
	"github.com/innotech/hydra/etcd"
	"github.com/innotech/hydra/server"
)

func main() {
	// Load configuration.
	var config = config.New()
	if err := config.Load(os.Args[1:]); err != nil {
		// fmt.Println(server.Usage() + "\n")
		fmt.Println(err.Error() + "\n")
		os.Exit(1)
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

	etcdDriver := driver.NewEtcdDriver(etcd.EtcdServer, etcd.PeerServer)
	var server = server.NewServer(etcdDriver)
	server.Start()
}
