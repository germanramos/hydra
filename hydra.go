package main

import (
	"github.com/innotech/hydra/config"
	"github.com/innotech/hydra/etcd"
)

func main() {
	// Load configuration.
	var config = config.New()
	var etcd = etcd.New()
	etcd.Start(config.EtcdConf)
}
