package config_test

import (
	. "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	// "fmt"
	"io/ioutil"
	"os"
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
				DATA_DIR     string = "/tmp/hydra-0"
				NAME         string = "hydra-0"
				PRIVATE_ADDR string = "127.0.0.1:8771"
				PUBLIC_ADDR  string = "127.0.0.1:8772"
			)
			fileContent := `
				data_dir = "` + DATA_DIR + `"
				name = "` + NAME + `"
				private_addr = "` + PRIVATE_ADDR + `"
				public_addr = "` + PUBLIC_ADDR + `"
			`
			WithTempFile(fileContent, func(pathToFile string) {
				c := New()
				err := c.LoadFile(pathToFile)
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
					Expect(c.DataDir).To(Equal(DATA_DIR))
					Expect(c.Name).To(Equal(NAME))
					Expect(c.PrivateAddr).To(Equal(PRIVATE_ADDR))
					Expect(c.PublicAddr).To(Equal(PUBLIC_ADDR))
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
			It("should be loaded successfully", func() {
				Expect(c.ConfigFilePath).To(Equal(DEFAULT_CONFIG_FILE_PATH))
				c.ConfigFilePath = "/no-data-config"
				err := c.Load([]string{})
				Expect(err).To(BeNil(), "error should be nil")
				Expect(c.DataDir).To(Equal(DEFAULT_DATA_DIR))
				Expect(c.PeerAddr).To(Equal(DEFAULT_PEER_ADDR))
				Expect(c.PrivateAddr).To(Equal(DEFAULT_PRIVATE_ADDR))
				Expect(c.PublicAddr).To(Equal(DEFAULT_PUBLIC_ADDR))
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
		Context("When -data-dir flag exists", func() {
			const DATA_DIR string = "/tmp/flag/"
			c := New()
			err := c.LoadFlags([]string{"-data-dir", DATA_DIR})
			It("should be loaded successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(c.DataDir).To(Equal(DATA_DIR))
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
		c := New()
		c.DataDir = "/tmp/hydra-tests"
		err := c.LoadEtcdConfig()
		It("should be loaded successfully", func() {
			Expect(err).To(BeNil(), "error should be nil")
		})
		It("should be override the default configuration loaded from default system configuration file", func() {
			Expect(c.EtcdConf.DataDir).To(Equal(c.DataDir), "c.EtcdConfig.DataDir should be equal "+c.DataDir)
		})
	})
})
