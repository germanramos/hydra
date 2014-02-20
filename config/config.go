package config

import (
	"flag"
	"io/ioutil"
	// "fmt"
	"os"

	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
	// "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
)

const (
	DefaultConfigFilePath = "/etc/hydra/hydra.conf"
	DEFAULT_ADDR          = "127.0.0.1:4001"
	DEFAULT_PEER_ADDR     = "127.0.0.1:7001"
)

// const DefaultConfigFilePath = "/etc/hydra/hydra.conf"

type Config struct {
	ConfigFilePath string
	EtcdConf       *etcdConfig.Config
	Addr           string
	PeerAddr       string
	// EtcdConf *config.Config
}

func New() *Config {
	c := new(Config)
	// conf.EtcdConf = new(etcdConfig.Config)
	c.EtcdConf = etcdConfig.New()
	c.ConfigFilePath = DefaultConfigFilePath
	c.PeerAddr = DEFAULT_PEER_ADDR

	return c
}

// Loads the configuration from the system config, command line config,
// environment variables, and finally command line arguments.
func (c *Config) Load(arguments []string) error {
	var path string
	f := flag.NewFlagSet("hydra", -1)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&path, "config", "", "path to config file")
	f.Parse(arguments)

	// Load from config file.
	if path != "" {
		// Load from config file specified in arguments.
		if err := c.LoadFile(path); err != nil {
			return err
		}
	} else {
		// Load from system file.
		if err := c.LoadSystemFile(); err != nil {
			return err
		}
	}

	// Load from command line flags.
	// if err := c.LoadFlags(arguments); err != nil {
	// 	return err
	// }

	// Load etcd configuration.
	// TODO: Fix -> bad flag returns an error
	if err := c.EtcdConf.Load(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

// Loads from the default hydra configuration file path if it exists.
func (c *Config) LoadSystemFile() error {
	if _, err := os.Stat(c.ConfigFilePath); os.IsNotExist(err) {
		return nil
	}
	return c.LoadFile(c.ConfigFilePath)
}

// Loads configuration from a file.
func (c *Config) LoadFile(path string) error {
	_, err := toml.DecodeFile(path, &c)
	return err
}

// Loads configuration from command line flags.
func (c *Config) LoadFlags(arguments []string) error {
	var ignoredString string

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	return nil
}

func (c *Config) LoadEtcdConfig() error {
	return nil
}
