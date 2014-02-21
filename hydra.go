package main

import (
	"fmt"
	"os"

	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/etcd"
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
		os.Exit(1)
	}

	// Load etcd configuration.
	if err := config.LoadEtcdConfig(); err != nil {
		os.Exit(1)
	}

	var etcd = etcd.New(config.EtcdConf)
	etcd.Start()
}
