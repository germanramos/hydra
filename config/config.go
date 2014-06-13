package config

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/innotech/hydra/log"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
	etcdConfig "github.com/innotech/hydra/vendors/github.com/coreos/etcd/config"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/discovery"
)

const (
	DEFAULT_APPS_FILE                = "/etc/hydra/apps.json"
	DEFAULT_BALANCE_TIMEOUT          = 2500
	DEFAULT_CONFIG_FILE_PATH         = "/etc/hydra/hydra.conf"
	DEFAULT_DATA_DIR                 = "/tmp/hydra"
	DEFAULT_ETCD_ADDR                = "127.0.0.1:7401"
	DEFAULT_INSTANCE_EXPIRATION_TIME = 30
	DEFAULT_LOAD_BALANCER_ADDR       = "*:7777"
	DEFAULT_PEER_ADDR                = "127.0.0.1:7701"
	DEFAULT_PEER_HEARTBEAT_TIMEOUT   = 50
	DEFAULT_PEER_ELECTION_TIMEOUT    = 200
	DEFAULT_PRIVATE_ADDR             = "127.0.0.1:7771"
	DEFAULT_PUBLIC_ADDR              = "127.0.0.1:7772"
	DEFAULT_SNAPSHOT                 = true
	DEFAULT_SNAPSHOT_COUNT           = 20000
	DEFAULT_VERBOSE                  = false
)

type Config struct {
	EtcdConf *etcdConfig.Config
	AppsConf *ApplicationsConfig

	AppsFile               string `toml:"apps_file"`
	BalanceTimeout         int    `toml:"balance_timeout"`
	BindAddr               string `toml:"bind_addr"`
	CAFile                 string `toml:"ca_file"`
	CertFile               string `toml:"cert_file"`
	ConfigFilePath         string
	DataDir                string `toml:"data_dir"`
	Discovery              string
	EtcdAddr               string `toml:"addr"`
	Force                  bool
	InstanceExpirationTime int    `toml:"instance_expiration_time"`
	KeyFile                string `toml:"key_file"`
	LoadBalancerAddr       string `toml:"load_balancer_addr"`
	Name                   string
	Peers                  []string
	PrivateAddr            string `toml:"private_addr"`
	PublicAddr             string `toml:"public_addr"`
	Snapshot               bool
	SnapshotCount          int  `toml:"snapshot_count"`
	Verbose                bool `toml:"verbose"`
	Peer                   struct {
		Addr             string `toml:"addr"`
		BindAddr         string `toml:"bind_addr"`
		CAFile           string `toml:"ca_file"`
		CertFile         string `toml:"cert_file"`
		KeyFile          string `toml:"key_file"`
		HeartbeatTimeout int    `toml:"heartbeat_timeout"`
		ElectionTimeout  int    `toml:"election_timeout"`
	}
}

func New() *Config {
	c := new(Config)
	c.EtcdConf = etcdConfig.New()
	c.AppsFile = DEFAULT_APPS_FILE
	c.BalanceTimeout = DEFAULT_BALANCE_TIMEOUT
	c.ConfigFilePath = DEFAULT_CONFIG_FILE_PATH
	c.DataDir = DEFAULT_DATA_DIR
	// c.EtcdAddr = DEFAULT_ETCD_ADDR
	c.InstanceExpirationTime = DEFAULT_INSTANCE_EXPIRATION_TIME
	c.LoadBalancerAddr = DEFAULT_LOAD_BALANCER_ADDR
	c.PrivateAddr = DEFAULT_PRIVATE_ADDR
	c.PublicAddr = DEFAULT_PUBLIC_ADDR
	c.Snapshot = DEFAULT_SNAPSHOT
	c.SnapshotCount = DEFAULT_SNAPSHOT_COUNT
	c.Verbose = DEFAULT_VERBOSE
	c.Peer.Addr = DEFAULT_PEER_ADDR
	c.Peer.HeartbeatTimeout = DEFAULT_PEER_HEARTBEAT_TIMEOUT
	c.Peer.ElectionTimeout = DEFAULT_PEER_ELECTION_TIMEOUT

	return c
}

// TODO: comment
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

	// Attempt cluster discovery
	if c.Discovery != "" {
		if err := c.handleDiscovery(); err != nil {
			return err
		}
	}

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
	var peers, ignoredString string

	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&c.EtcdAddr, "addr", c.EtcdAddr, "")
	f.StringVar(&c.AppsFile, "apps-file", c.AppsFile, "")
	f.IntVar(&c.BalanceTimeout, "balance-timeout", c.BalanceTimeout, "")
	f.StringVar(&c.BindAddr, "bind-addr", c.BindAddr, "")
	f.StringVar(&c.CAFile, "ca-file", c.CAFile, "")
	f.StringVar(&c.CertFile, "cert-file", c.CertFile, "")
	f.StringVar(&c.DataDir, "data-dir", c.DataDir, "")
	f.StringVar(&c.Discovery, "discovery", c.Discovery, "")
	f.IntVar(&c.InstanceExpirationTime, "instance-expiration-time", c.InstanceExpirationTime, "")
	f.StringVar(&c.KeyFile, "key-file", c.KeyFile, "")
	f.BoolVar(&c.Force, "f", false, "")
	f.BoolVar(&c.Force, "force", false, "")
	f.StringVar(&c.LoadBalancerAddr, "load-balancer-addr", c.LoadBalancerAddr, "")
	f.StringVar(&c.Name, "name", c.Name, "")
	f.StringVar(&peers, "peers", "", "")
	f.StringVar(&c.PrivateAddr, "private-addr", c.PrivateAddr, "")
	f.StringVar(&c.PublicAddr, "public-addr", c.PublicAddr, "")
	f.BoolVar(&c.Snapshot, "snapshot", true, "")
	f.IntVar(&c.SnapshotCount, "snapshot-count", c.SnapshotCount, "")
	f.BoolVar(&c.Verbose, "v", c.Verbose, "")
	f.BoolVar(&c.Verbose, "verbose", c.Verbose, "")

	f.StringVar(&c.Peer.Addr, "peer-addr", c.Peer.Addr, "")
	f.StringVar(&c.Peer.BindAddr, "peer-bind-addr", c.Peer.BindAddr, "")
	f.StringVar(&c.Peer.CAFile, "peer-ca-file", c.Peer.CAFile, "")
	f.StringVar(&c.Peer.CertFile, "peer-cert-file", c.Peer.CertFile, "")
	f.StringVar(&c.Peer.KeyFile, "peer-key-file", c.Peer.KeyFile, "")
	f.IntVar(&c.Peer.HeartbeatTimeout, "peer-heartbeat-timeout", c.Peer.HeartbeatTimeout, "")
	f.IntVar(&c.Peer.ElectionTimeout, "peer-election-timeout", c.Peer.ElectionTimeout, "")

	// BEGIN IGNORED FLAGS
	f.StringVar(&ignoredString, "config", "", "")
	// BEGIN IGNORED FLAGS

	if err := f.Parse(arguments); err != nil {
		return err
	}

	// Convert some parameters to lists.
	if peers != "" {
		c.Peers = strings.Split(peers, ",")
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
	if c.EtcdAddr == "" {
		addLineToFileContent(`addr = "` + DEFAULT_ETCD_ADDR + `"`)
	} else {
		addLineToFileContent(`addr = "` + c.EtcdAddr + `"`)
	}
	addLineToFileContent(`ca_file = "` + c.CAFile + `"`)
	addLineToFileContent(`cert_file = "` + c.CertFile + `"`)
	addLineToFileContent(`data_dir = "` + c.DataDir + `"`)
	addLineToFileContent(`key_file = "` + c.KeyFile + `"`)
	addLineToFileContent(`name = "` + c.Name + `"`)
	peers := ""
	for i, addr := range c.Peers {
		if i > 0 {
			peers = peers + ", "
		}
		peers = peers + `"` + addr + `"`
	}
	addLineToFileContent(`peers = [` + peers + `]`)
	addLineToFileContent(`snapshot = ` + strconv.FormatBool(c.Snapshot))
	addLineToFileContent(`snapshot_count = ` + strconv.FormatInt(int64(c.SnapshotCount), 10))
	addLineToFileContent(`[peer]`)
	addLineToFileContent(`addr = "` + c.Peer.Addr + `"`)
	addLineToFileContent(`bind_addr = "` + c.Peer.BindAddr + `"`)
	addLineToFileContent(`ca_file = "` + c.Peer.CAFile + `"`)
	addLineToFileContent(`cert_file = "` + c.Peer.CertFile + `"`)
	addLineToFileContent(`key_file = "` + c.Peer.KeyFile + `"`)
	addLineToFileContent(`heartbeat_timeout = ` + strconv.FormatInt(int64(c.Peer.HeartbeatTimeout), 10))
	addLineToFileContent(`election_timeout = ` + strconv.FormatInt(int64(c.Peer.ElectionTimeout), 10))

	return content
}

func (c *Config) handleDiscovery() error {
	p, err := discovery.Do(c.Discovery, c.Name, c.Peer.Addr)

	// This is fatal, discovery encountered an unexpected error
	// and we have no peer list.
	if err != nil && len(c.Peers) == 0 {
		log.Fatalf("Discovery failed and a backup peer list wasn't provided: %v", err)
		return err
	}

	// Warn about errors coming from discovery, this isn't fatal
	// since the user might have provided a peer list elsewhere.
	if err != nil {
		log.Warnf("Discovery encountered an error but a backup peer list (%v) was provided: %v", c.Peers, err)
	}

	for i := range p {
		// Strip the scheme off of the peer if it has one
		// TODO(bp): clean this up!
		purl, err := url.Parse(p[i])
		if err == nil {
			p[i] = purl.Host
		}
	}

	c.Peers = p

	return nil
}
