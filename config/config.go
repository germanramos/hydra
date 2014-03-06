package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
)

const (
	DEFAULT_CONFIG_FILE_PATH = "/etc/hydra/hydra.conf"
	DEFAULT_DATA_DIR         = "./"
	DEFAULT_PEER_ADDR        = "127.0.0.1:7001"
	DEFAULT_PRIVATE_ADDR     = "127.0.0.1:7771"
	DEFAULT_PUBLIC_ADDR      = "127.0.0.1:7772"
)

type Config struct {
	EtcdConf *etcdConfig.Config

	ConfigFilePath string
	DataDir        string `toml:"data_dir"`
	Force          bool
	Name           string
	Peers          []string
	PeerAddr       string `toml:"peer_addr"`
	PrivateAddr    string `toml:"private_addr"`
	PublicAddr     string `toml:"public_addr"`
}

func New() *Config {
	c := new(Config)
	c.EtcdConf = etcdConfig.New()
	c.ConfigFilePath = DEFAULT_CONFIG_FILE_PATH
	c.DataDir = DEFAULT_DATA_DIR
	c.PeerAddr = DEFAULT_PEER_ADDR
	c.PrivateAddr = DEFAULT_PRIVATE_ADDR
	c.PublicAddr = DEFAULT_PUBLIC_ADDR

	return c
}

// Loads the configuration from the system config, command line config,
// environment variables, and finally command line arguments.
func (c *Config) Load(arguments []string) error {
	var path string
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

	// TODO: name is required make default or check if exist

	// Force remove server configuration if specified.
	if c.Force {
		if err := c.Reset(); err != nil {
			return err
		}
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
	f.StringVar(&c.DataDir, "data-dir", c.DataDir, "")
	f.BoolVar(&c.Force, "f", false, "")
	f.BoolVar(&c.Force, "force", false, "")
	f.StringVar(&c.Name, "name", c.Name, "")
	f.StringVar(&c.PrivateAddr, "private-addr", c.PrivateAddr, "")
	f.StringVar(&c.PublicAddr, "public-addr", c.PublicAddr, "")
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
	defer os.Remove(f.Name())
	if err := c.EtcdConf.Load([]string{"-config", f.Name()}); err != nil {
		return err
	}

	return nil
}

// Reset removes all server configuration files.
func (c *Config) Reset() error {
	if err := os.RemoveAll(filepath.Join(c.DataDir, "log")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(c.DataDir, "conf")); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(c.DataDir, "snapshot")); err != nil {
		return err
	}

	return nil
}

func (c *Config) makeEtcdConfig() string {
	var content string
	addLineToFileContent := func( /*fileContent *string, */ line string) {
		// *fileContent = *fileContent + line + "\n"
		content = content + line + "\n"
	}
	// addLineToFileContent(`addr = "` + c.Addr + `"`)
	addLineToFileContent(`data_dir = "` + c.DataDir + `"`)
	addLineToFileContent(`name = "` + c.Name + `"`)
	peers := ""
	for i, addr := range c.Peers {
		if i > 0 {
			peers = peers + ", "
		}
		peers = peers + `"` + addr + `"`
	}
	addLineToFileContent(`peers = [` + peers + `]`)
	addLineToFileContent(`[peer]`)
	addLineToFileContent(`addr = "` + c.PeerAddr + `"`)

	return content
}
