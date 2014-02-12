package config

import (
	"fmt"
	"os"

	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
)

const DefaultConfigFilePath = "/etc/hydra/hydra.conf"

type Config struct {
	EtcdConf       *etcdConfig.Config
	ConfigFilePath string

	Addr string `toml:"addr"`
}

func New() *Config {
	conf := new(Config)
	conf.EtcdConf = new(etcdConfig.Config)
	conf.ConfigFilePath = DefaultConfigFilePath
	conf.Addr = "127.0.0.1:4001"

	if _, err := os.Stat(conf.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}

	if err := conf.LoadHydraConfig(); err != nil {
		// TODO: log?
		fmt.Println(err.Error() + "\n")
		os.Exit(1)
	}

	if err := conf.LoadEtcdConfig(); err != nil {
		// TODO: log?
		fmt.Println(err.Error() + "\n")
		os.Exit(1)
	}
	return conf
}

func (conf *Config) LoadHydraConfig() error {
	_, err := toml.DecodeFile(conf.ConfigFilePath, &conf)
	return err
}

func (conf *Config) LoadEtcdConfig() error {
	_, err := toml.DecodeFile(conf.ConfigFilePath, &conf.EtcdConf)
	return err
}

// func (conf *Config) Load() error {
// 	if err := conf.LoadCofigFile(conf.ConfigFilePath); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Load configuration from the configuration file
// func (conf *Config) LoadCofigFile(filePath string, config *) error {
// 	_, err := toml.Decode(filePath, &conf)
// 	return err
// }
