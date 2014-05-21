package config_test

import (
	. "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	// "fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var _ = Describe("Config", func() {
	// HELPERS ////////////////////////////////////////////////////////////////////////
	WithTempFile := func(content string, fn func(string)) {
		f, _ := ioutil.TempFile("", "")
		f.WriteString(content)
		f.Close()
		defer os.Remove(f.Name())
		fn(f.Name())
	}
	// END OF HELPERS /////////////////////////////////////////////////////////////////
	Describe("loading from TOML", func() {
		Context("when the TOML file exists", func() {
			const (
				APPS_FILE          string = "/etc/hydra-test/apps.json"
				BIND_ADDR          string = "127.0.0.1:5011"
				CA_FILE            string = "./fixtures/ca/server-chain.pem"
				CERT_FILE          string = "./fixtures/ca/server.crt"
				DATA_DIR           string = "/tmp/hydra-0"
				DISCOVERY          string = "http://etcd.local:4001/v2/keys/_etcd/registry/examplecluster"
				ETCD_ADDR          string = "127.0.0.1:5001"
				KEY_FILE           string = "./fixtures/ca/server.key.insecure"
				LOAD_BALANCER_ADDR string = "*:5555"
				NAME               string = "hydra-0"
				PEER_1             string = "192.168.113.101:7001"
				PEER_2             string = "192.168.113.102:7001"
				PRIVATE_ADDR       string = "127.0.0.1:8771"
				PUBLIC_ADDR        string = "127.0.0.1:8772"
				SNAPSHOT           bool   = false
				SNAPSHOT_COUNT     int    = 333

				PEER_ADDR              string = "127.0.0.1:8001"
				PEER_BIND_ADDR         string = "127.0.0.1:8011"
				PEER_CA_FILE           string = "./fixtures/ca/peer_server-chain.pem"
				PEER_CERT_FILE         string = "./fixtures/ca/peer_server.crt"
				PEER_KEY_FILE          string = "./fixtures/ca/peer_server.key.insecure"
				PEER_HEARTBEAT_TIMEOUT int    = 55
				PEER_ELECTION_TIMEOUT  int    = 207
			)
			fileContent := `
				addr = "` + ETCD_ADDR + `"
				bind_addr = "` + BIND_ADDR + `"
				apps_file = "` + APPS_FILE + `"
				ca_file = "` + CA_FILE + `"
				cert_file = "` + CERT_FILE + `"
				data_dir = "` + DATA_DIR + `"
				discovery = "` + DISCOVERY + `"
				key_file = "` + KEY_FILE + `"
				load_balancer_addr = "` + LOAD_BALANCER_ADDR + `"
				name = "` + NAME + `"
				peers = ["` + PEER_1 + `","` + PEER_2 + `"]
				private_addr = "` + PRIVATE_ADDR + `"
				public_addr = "` + PUBLIC_ADDR + `"
				snapshot = ` + strconv.FormatBool(SNAPSHOT) + `
				snapshot_count = ` + strconv.FormatInt(int64(SNAPSHOT_COUNT), 10) + `
				[peer]
				addr = "` + PEER_ADDR + `"
				bind_addr = "` + PEER_BIND_ADDR + `"
				ca_file = "` + PEER_CA_FILE + `"
				cert_file = "` + PEER_CERT_FILE + `"
				key_file = "` + PEER_KEY_FILE + `"
				heartbeat_timeout = ` + strconv.FormatInt(int64(PEER_HEARTBEAT_TIMEOUT), 10) + `
				election_timeout = ` + strconv.FormatInt(int64(PEER_ELECTION_TIMEOUT), 10) + ` 
			`
			WithTempFile(fileContent, func(pathToFile string) {
				c := New()
				err := c.LoadFile(pathToFile)
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
					Expect(c.AppsFile).To(Equal(APPS_FILE))
					Expect(c.BindAddr).To(Equal(BIND_ADDR))
					Expect(c.CAFile).To(Equal(CA_FILE))
					Expect(c.CertFile).To(Equal(CERT_FILE))
					Expect(c.DataDir).To(Equal(DATA_DIR))
					Expect(c.Discovery).To(Equal(DISCOVERY))
					Expect(c.EtcdAddr).To(Equal(ETCD_ADDR))
					Expect(c.KeyFile).To(Equal(KEY_FILE))
					Expect(c.LoadBalancerAddr).To(Equal(LOAD_BALANCER_ADDR))
					Expect(c.Name).To(Equal(NAME))
					Expect(c.Peers).To(HaveLen(2))
					Expect(c.Peers).To(ContainElement(PEER_1))
					Expect(c.Peers).To(ContainElement(PEER_2))
					Expect(c.PrivateAddr).To(Equal(PRIVATE_ADDR))
					Expect(c.PublicAddr).To(Equal(PUBLIC_ADDR))
					Expect(c.Snapshot).To(Equal(SNAPSHOT))
					Expect(c.SnapshotCount).To(Equal(SNAPSHOT_COUNT))
					Expect(c.Peer.Addr).To(Equal(PEER_ADDR))
					Expect(c.Peer.BindAddr).To(Equal(PEER_BIND_ADDR))
					Expect(c.Peer.CAFile).To(Equal(PEER_CA_FILE))
					Expect(c.Peer.CertFile).To(Equal(PEER_CERT_FILE))
					Expect(c.Peer.KeyFile).To(Equal(PEER_KEY_FILE))
					Expect(c.Peer.HeartbeatTimeout).To(Equal(PEER_HEARTBEAT_TIMEOUT))
					Expect(c.Peer.ElectionTimeout).To(Equal(PEER_ELECTION_TIMEOUT))
				})
			})
		})
	})

	Describe("loading from command flags", func() {
		Context("when config flag exists", func() {
			const FLAG_VALUE string = "/etc/hydra/hydra.conf"
			c := New()
			c.LoadFlags([]string{"-config", FLAG_VALUE})
			It("should be loaded successfully", func() {
				Expect(c.ConfigFilePath).To(Equal(FLAG_VALUE))
			})
		})
	})

	Describe("loading without flags", func() {
		Context("when default system cofig file doesn't exist", func() {
			c := New()
			It("should be loaded the default configuration successfully", func() {
				Expect(c.ConfigFilePath).To(Equal(DEFAULT_CONFIG_FILE_PATH))
				c.ConfigFilePath = "/no-data-config"
				err := c.Load([]string{})
				Expect(err).To(BeNil(), "error should be nil")
				Expect(c.AppsFile).To(Equal(DEFAULT_APPS_FILE))
				Expect(c.DataDir).To(Equal(DEFAULT_DATA_DIR))
				// Expect(c.EtcdAddr).To(Equal(DEFAULT_ETCD_ADDR))
				Expect(c.LoadBalancerAddr).To(Equal(DEFAULT_LOAD_BALANCER_ADDR))
				Expect(c.Peer.Addr).To(Equal(DEFAULT_PEER_ADDR))
				Expect(c.Peer.HeartbeatTimeout).To(Equal(DEFAULT_PEER_HEARTBEAT_TIMEOUT))
				Expect(c.Peer.ElectionTimeout).To(Equal(DEFAULT_PEER_ELECTION_TIMEOUT))
				Expect(c.PrivateAddr).To(Equal(DEFAULT_PRIVATE_ADDR))
				Expect(c.PublicAddr).To(Equal(DEFAULT_PUBLIC_ADDR))
				Expect(c.Snapshot).To(Equal(DEFAULT_SNAPSHOT))
				Expect(c.SnapshotCount).To(Equal(DEFAULT_SNAPSHOT_COUNT))
			})
		})
		Context("when default system cofig file exists", func() {
			systemPublicAddr := DEFAULT_PUBLIC_ADDR + "0"
			systemFileContent := `public_addr = "` + systemPublicAddr + `"`
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				err := c.Load([]string{})
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
				It("should be override the default configuration", func() {
					Expect(c.PublicAddr).To(Equal(systemPublicAddr), "c.Addr should be equal "+systemPublicAddr)
				})
			})
		})
	})

	Describe("loading from flags", func() {
		Context("when bad flag exists", func() {
			c := New()
			err := c.LoadFlags([]string{"-bad-flag"})
			It("should be throw an error", func() {
				Expect(err).To(HaveOccurred(), "No bad flag are allowed")
			})
			It("should be have an specific error message", func() {
				Expect(err.Error()).To(Equal("flag provided but not defined: -bad-flag"))
			})
		})
		Context("When -addr flag exists", func() {
			const ETCD_ADDR string = "127.0.0.1:6001"
			c := New()
			err := c.LoadFlags([]string{"-addr", ETCD_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.EtcdAddr).To(Equal(ETCD_ADDR))
			})
		})
		Context("When -bind-addr flag exists", func() {
			const BIND_ADDR string = "127.0.0.1:6011"
			c := New()
			err := c.LoadFlags([]string{"-bind-addr", BIND_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.BindAddr).To(Equal(BIND_ADDR))
			})
		})
		Context("When -apps-file flag exists", func() {
			const APPS_FILE string = "/tmp/apps.json"
			c := New()
			err := c.LoadFlags([]string{"-apps-file", APPS_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.AppsFile).To(Equal(APPS_FILE))
			})
		})
		Context("When -ca-file flag exists", func() {
			const CA_FILE string = "./fixtures/ca/server-chain_1.pem"
			c := New()
			err := c.LoadFlags([]string{"-ca-file", CA_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.CAFile).To(Equal(CA_FILE))
			})
		})
		Context("When -cert-file flag exists", func() {
			const CERT_FILE string = "./fixtures/ca/server_1.crt"
			c := New()
			err := c.LoadFlags([]string{"-cert-file", CERT_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.CertFile).To(Equal(CERT_FILE))
			})
		})
		Context("When -data-dir flag exists", func() {
			const DATA_DIR string = "/tmp/flag/"
			c := New()
			err := c.LoadFlags([]string{"-data-dir", DATA_DIR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.DataDir).To(Equal(DATA_DIR))
			})
		})
		Context("When -discovery flag exists", func() {
			const DISCOVERY string = "http://etcd.local:4001/v2/keys/_etcd/registry/examplecluster_1"
			c := New()
			err := c.LoadFlags([]string{"-discovery", DISCOVERY})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Discovery).To(Equal(DISCOVERY))
			})
		})
		Context("When -discovery flag exists", func() {
			const DISCOVERY string = "http://etcd.local:4001/v2/keys/_etcd/registry/examplecluster_1"
			c := New()
			err := c.LoadFlags([]string{"-discovery", DISCOVERY})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Discovery).To(Equal(DISCOVERY))
			})
		})
		Context("When -f flag exists", func() {
			c := New()
			err := c.LoadFlags([]string{"-f"})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Force).To(BeTrue())
			})
		})
		Context("When -force flag exists", func() {
			c := New()
			err := c.LoadFlags([]string{"-force"})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Force).To(BeTrue())
			})
		})
		Context("When -key-file flag exists", func() {
			const KEY_FILE string = "./fixtures/ca/server_1.key.insecure"
			c := New()
			err := c.LoadFlags([]string{"-key-file", KEY_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.KeyFile).To(Equal(KEY_FILE))
			})
		})
		Context("When -load-balancer-addr flag exists", func() {
			const LOAD_BALANCER_ADDR string = "*:2222"
			c := New()
			err := c.LoadFlags([]string{"-load-balancer-addr", LOAD_BALANCER_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.LoadBalancerAddr).To(Equal(LOAD_BALANCER_ADDR))
			})
		})
		Context("When -peer-addr flag exists", func() {
			const PEER_ADDR string = "127.0.0.1:8001"
			c := New()
			err := c.LoadFlags([]string{"-peer-addr", PEER_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.Addr).To(Equal(PEER_ADDR))
			})
		})
		Context("When -peer-bind-addr flag exists", func() {
			const PEER_BIND_ADDR string = "127.0.0.1:8011"
			c := New()
			err := c.LoadFlags([]string{"-peer-bind-addr", PEER_BIND_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.BindAddr).To(Equal(PEER_BIND_ADDR))
			})
		})
		Context("When -peer-ca-file flag exists", func() {
			const PEER_CA_FILE string = "./fixtures/ca/peer_server-chain.pem"
			c := New()
			err := c.LoadFlags([]string{"-peer-ca-file", PEER_CA_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.CAFile).To(Equal(PEER_CA_FILE))
			})
		})
		Context("When -peer-heartbeat-timeout flag exists", func() {
			const PEER_CERT_FILE string = "./fixtures/ca/peer_server.crt"
			c := New()
			err := c.LoadFlags([]string{"-peer-cert-file", PEER_CERT_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.CertFile).To(Equal(PEER_CERT_FILE))
			})
		})
		Context("When -peer-cert-file flag exists", func() {
			const PEER_KEY_FILE string = "./fixtures/ca/peer_server.key.insecure"
			c := New()
			err := c.LoadFlags([]string{"-peer-key-file", PEER_KEY_FILE})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.KeyFile).To(Equal(PEER_KEY_FILE))
			})
		})
		Context("When -peer-heartbeat-timeout flag exists", func() {
			const PEER_HEARTBEAT_TIMEOUT int = 21
			c := New()
			err := c.LoadFlags([]string{"-peer-heartbeat-timeout", strconv.FormatInt(int64(PEER_HEARTBEAT_TIMEOUT), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.HeartbeatTimeout).To(Equal(PEER_HEARTBEAT_TIMEOUT))
			})
		})
		Context("When -peer-election-timeout flag exists", func() {
			const PEER_ELECTION_TIMEOUT int = 201
			c := New()
			err := c.LoadFlags([]string{"-peer-election-timeout", strconv.FormatInt(int64(PEER_ELECTION_TIMEOUT), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peer.ElectionTimeout).To(Equal(PEER_ELECTION_TIMEOUT))
			})
		})
		Context("When -peers flag exists", func() {
			const PEER_1 string = "203.0.113.101:7001"
			const PEER_2 string = "203.0.113.102:7001"
			c := New()
			err := c.LoadFlags([]string{"-peers", PEER_1 + "," + PEER_2})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Peers).To(HaveLen(2))
				Expect(c.Peers).To(ContainElement(PEER_1))
				Expect(c.Peers).To(ContainElement(PEER_2))
			})
		})
		Context("When -private-addr flag exists", func() {
			const PRIVATE_ADDR string = "localhost:4444"
			c := New()
			err := c.LoadFlags([]string{"-private-addr", PRIVATE_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.PrivateAddr).To(Equal(PRIVATE_ADDR))
			})
		})
		Context("When -public-addr flag exists", func() {
			const PUBLIC_ADDR string = "localhost:5555"
			c := New()
			err := c.LoadFlags([]string{"-public-addr", PUBLIC_ADDR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.PublicAddr).To(Equal(PUBLIC_ADDR))
			})
		})
		Context("When -snapshot flag exists", func() {
			c := New()
			c.Snapshot = false
			err := c.LoadFlags([]string{"-snapshot"})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Snapshot).To(BeTrue())
			})
		})
		Context("When -snapshot-count flag exists", func() {
			const SNAPSHOT_COUNT int = 777
			c := New()
			err := c.LoadFlags([]string{"-snapshot-count", strconv.FormatInt(int64(SNAPSHOT_COUNT), 10)})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.SnapshotCount).To(Equal(SNAPSHOT_COUNT))
			})
		})
		Context("When -name flag exists", func() {
			const NAME string = "test-0"
			c := New()
			err := c.LoadFlags([]string{"-name", NAME})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.Name).To(Equal(NAME))
			})
		})
		Context("when config flag exists", func() {
			Context("and no more flags exist", func() {
				systemPublicAddr := DEFAULT_PUBLIC_ADDR + "0"
				systemFileContent := `public_addr = "` + systemPublicAddr + `"`
				customPublicAddr := systemPublicAddr + "0"
				customFileContent := `public_addr = "` + customPublicAddr + `"`
				WithTempFile(systemFileContent, func(pathToSystemFile string) {
					WithTempFile(customFileContent, func(pathToCustomFile string) {
						c := New()
						c.ConfigFilePath = pathToSystemFile
						err := c.Load([]string{"-config", pathToCustomFile})
						It("should be loaded successfully", func() {
							Expect(err).To(BeNil(), "error should be nil")
							Expect(c.ConfigFilePath).To(Equal(pathToSystemFile))
						})
						It("should be override the default configuration loaded from default system configuration file", func() {
							Expect(c.PublicAddr).To(Equal(customPublicAddr), "c.Addr should be equal "+customPublicAddr)
						})
					})
				})
			})
			Context("and also more valid flags exist", func() {
				customPublicAddr := DEFAULT_PUBLIC_ADDR + "0"
				customFileContent := `public_addr = "` + customPublicAddr + `"`
				addrCustomFlag := customPublicAddr + "0"
				WithTempFile(customFileContent, func(pathToCustomFile string) {
					c := New()
					err := c.Load([]string{"-public-addr", addrCustomFlag, "-config", pathToCustomFile})
					It("should be loaded successfully", func() {
						Expect(err).To(BeNil(), "error should be nil")
					})
					It("should be override the configuration loaded from custom configuration file", func() {
						Expect(c.PublicAddr).To(Equal(addrCustomFlag), "c.Addr should be equal "+addrCustomFlag)
					})
					// TODO: Check Error() message
				})
			})
		})
		Context("when default system cofig file doesn't exist", func() {
			systemPublicAddr := DEFAULT_PUBLIC_ADDR + "0"
			systemFileContent := `public_addr = "` + systemPublicAddr + `"`
			customPublicAddr := systemPublicAddr + "0"
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				err := c.Load([]string{"-public-addr", customPublicAddr})
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
				It("should be override the default configuration loaded from default system configuration file", func() {
					Expect(c.PublicAddr).To(Equal(customPublicAddr), "c.Addr should be equal "+customPublicAddr)
				})
			})
		})
	})

	Describe("loading etcd configuration", func() {
		Context("When transport Security with HTTPS will NOT be enabled", func() {
			c := New()
			c.DataDir = "/tmp/hydra-tests_11"
			c.EtcdAddr = "127.0.0.1:4411"
			c.Snapshot = false
			c.SnapshotCount = 222
			c.Peer.Addr = "127.0.0.1:7711"
			c.Peer.BindAddr = "127.0.0.1:7722"
			c.Peer.HeartbeatTimeout = 23
			c.Peer.ElectionTimeout = 203
			err := c.LoadEtcdConfig()
			It("should be loaded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
			})
			It("should be override the default configuration loaded from default system configuration file", func() {
				Expect(c.EtcdConf.DataDir).To(Equal(c.DataDir))
				Expect(c.EtcdConf.Addr).To(Equal("http://" + c.EtcdAddr))
				Expect(c.EtcdConf.Snapshot).To(Equal(c.Snapshot))
				Expect(c.EtcdConf.SnapshotCount).To(Equal(c.SnapshotCount))
				Expect(c.EtcdConf.Peer.Addr).To(Equal("http://" + c.Peer.Addr))
				Expect(c.EtcdConf.Peer.BindAddr).To(Equal(c.Peer.BindAddr))
				Expect(c.EtcdConf.Peer.HeartbeatTimeout).To(Equal(c.Peer.HeartbeatTimeout))
				Expect(c.EtcdConf.Peer.ElectionTimeout).To(Equal(c.Peer.ElectionTimeout))
			})
		})
		Context("When transport Security with HTTPS will be enabled", func() {
			c := New()
			c.CAFile = "./fixtures/ca/server-chain.pem"
			c.CertFile = "./fixtures/ca/server.crt"
			c.EtcdAddr = "127.0.0.1:4411"
			c.KeyFile = "./fixtures/ca/server.key.insecure"
			c.Peer.Addr = "127.0.0.1:7711"
			c.Peer.CAFile = "./fixtures/ca/peer-server-chain.pem"
			c.Peer.CertFile = "./fixtures/ca/peerserver.crt"
			c.Peer.KeyFile = "./fixtures/ca/peerserver.key.insecure"
			err := c.LoadEtcdConfig()
			It("should be loaded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
			})
			It("should be override the default configuration loaded from default system configuration file", func() {
				Expect(c.EtcdConf.CAFile).To(Equal(c.CAFile))
				Expect(c.EtcdConf.CertFile).To(Equal(c.CertFile))
				Expect(c.EtcdConf.Addr).To(Equal("https://" + c.EtcdAddr))
				Expect(c.EtcdConf.KeyFile).To(Equal(c.KeyFile))
				Expect(c.EtcdConf.Peer.Addr).To(Equal("https://" + c.Peer.Addr))
				Expect(c.EtcdConf.Peer.CAFile).To(Equal(c.Peer.CAFile))
				Expect(c.EtcdConf.Peer.CertFile).To(Equal(c.Peer.CertFile))
				Expect(c.EtcdConf.Peer.KeyFile).To(Equal(c.Peer.KeyFile))
			})
		})

	})
})
