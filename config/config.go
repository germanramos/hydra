package config

import (
	"flag"
	// "fmt"
	// "errors"
	"io/ioutil"
	"os"
	// "strconv"

	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
	// "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
)

const (
	DefaultConfigFilePath = "/etc/hydra/hydra.conf"
	DEFAULT_ADDR          = "127.0.0.1:4001"
	DEFAULT_DATA_DIR      = "/tmp/hydra/"
	DEFAULT_PEER_ADDR     = "127.0.0.1:7001"
)

type Config struct {
	ConfigFilePath string
	EtcdConf       *etcdConfig.Config
	Addr           string
	DataDir        string
	PeerAddr       string
	// EtcdConf *config.Config
}

func New() *Config {
	c := new(Config)
	// conf.EtcdConf = new(etcdConfig.Config)
	c.EtcdConf = etcdConfig.New()
	c.ConfigFilePath = DefaultConfigFilePath
	c.Addr = DEFAULT_ADDR
	// c.DataDir = DEFAULT_DATA_DIR
	c.PeerAddr = DEFAULT_PEER_ADDR

	return c
}

// Loads the configuration from the system config, command line config,
// environment variables, and finally command line arguments.
func (c *Config) Load(arguments []string) error {
	var path string
	// f := flag.NewFlagSet("hydra", flag.ExitOnError)
	f := flag.NewFlagSet("hydra", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&path, "config", "", "path to config file")
	f.Parse(arguments)

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
	if err := c.LoadFlags(arguments); err != nil {
		return err
	}

	// if c.DataDir == "" {
	// 	// TODO: include log system
	// 	// log.Fatal("The data dir was not set and could not be guessed from machine name")
	// 	return errors.New("data directory attribute is required")
	// }

	// // Load etcd configuration.
	// if err := c.LoadEtcdConfig(); err != nil {
	// 	return err
	// }

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

	// f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&c.Addr, "addr", c.Addr, "")

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	return nil
}

// Loads etcd configuration
func (c *Config) LoadEtcdConfig() error {
	fileContent := c.makeEtcdConfig()
	// TODO: Check if file is created
	f, _ := ioutil.TempFile("", "")
	f.WriteString(fileContent)
	f.Close()
	// fmt.Println("etcd file: " + f.Name())
	defer os.Remove(f.Name())
	// c.WithTempFile(fileContent, func(pathToEtcdConfigFile string) {
	if err := c.EtcdConf.Load([]string{"-config", f.Name()}); err != nil {
		return err
	}

	return nil
	// })
}

// func (c *Config) WithTempFile(content string, fn func(string)) {
// 	f, _ := ioutil.TempFile("", "")
// 	f.WriteString(content)
// 	f.Close()
// 	fmt.Println("etcd file: " + f.Name())
// 	// defer os.Remove(f.Name())
// 	fn(f.Name())
// }

func (c *Config) makeEtcdConfig() string {
	var content string
	addLineToFileContent := func( /*fileContent *string, */ line string) {
		// *fileContent = *fileContent + line + "\n"
		content = content + line + "\n"
	}
	addLineToFileContent(`addr = "` + c.Addr + `"`)
	addLineToFileContent(`data_dir = "` + c.DataDir + `"`)

	return content
}
