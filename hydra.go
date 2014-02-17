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

	var etcd = etcd.New(config.EtcdConf)
	etcd.Start()
}
