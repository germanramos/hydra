package config_test

import (
	. "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"
)

var _ = Describe("Config", func() {
	WithTempFile := func(content string, fn func(string)) {
		f, _ := ioutil.TempFile("", "")
		f.WriteString(content)
		f.Close()
		defer os.Remove(f.Name())
		fn(f.Name())
	}

	Describe("loading from TOML", func() {
		Context("when the TOML file exists", func() {
			fileContent := `
				addr = "127.0.0.1:4002"
			`
			WithTempFile(fileContent, func(pathToFile string) {
				c := New()
				err := c.LoadFile(pathToFile)
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
			})
		})
	})

	Describe("loading from command flags", func() {
		Context("when config flag exists", func() {
			const _flagValue string = "/etc/hydra/hydra.conf"
			c := New()
			c.LoadFlags([]string{"-config", _flagValue})
			It("should be loaded successfully", func() {
				Expect(c.ConfigFilePath).To(Equal(_flagValue))
			})
		})
	})

	// Describe("loading with no data directory", func() {
	// 	c := New()
	// 	err := c.Load([]string{})
	// 	It("should be throw an error", func() {
	// 		Expect(err).To(HaveOccurred(), "data directory attribute should not be empty")
	// 	})
	// 	It("should be have an specific error message", func() {
	// 		Expect(err.Error()).To(Equal("data directory attribute is required"))
	// 	})
	// })

	Describe("loading without flags", func() {
		Context("when default system cofig file doesn't exist", func() {
			c := New()
			err := c.Load([]string{})
			It("should be loaded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
				Expect(c.Addr).To(Equal(DEFAULT_ADDR), "c.Addr should be equal "+DEFAULT_ADDR)
				// TODO: check all constants
			})
		})
		Context("when default system cofig file exists", func() {
			systemAddr := DEFAULT_ADDR + "0"
			systemFileContent := `addr = "` + systemAddr + `"`
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				err := c.Load([]string{})
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
				It("should be override the default configuration", func() {
					Expect(c.Addr).To(Equal(systemAddr), "c.Addr should be equal "+systemAddr)
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
		Context("when cofig flag exists", func() {
			Context("and no more flags exist", func() {
				systemAddr := DEFAULT_ADDR + "0"
				systemFileContent := `addr = "` + systemAddr + `"`
				customAddr := systemAddr + "0"
				customFileContent := `addr = "` + customAddr + `"`
				WithTempFile(systemFileContent, func(pathToSystemFile string) {
					WithTempFile(customFileContent, func(pathToCustomFile string) {
						c := New()
						c.ConfigFilePath = pathToSystemFile
						err := c.Load([]string{"-config", pathToCustomFile})
						It("should be loaded successfully", func() {
							Expect(err).To(BeNil(), "error should be nil")
						})
						It("should be override the default configuration loaded from default system configuration file", func() {
							Expect(c.Addr).To(Equal(customAddr), "c.Addr should be equal "+customAddr)
						})
					})
				})
			})
			Context("and also more valid flags exist", func() {
				customAddr := DEFAULT_ADDR + "0"
				customFileContent := `addr = "` + customAddr + `"`
				addrCustomFlag := customAddr + "0"
				WithTempFile(customFileContent, func(pathToCustomFile string) {
					c := New()
					fmt.Println("&&&&&&&&&& ")
					err := c.Load([]string{"-addr", addrCustomFlag, "-config", pathToCustomFile})
					It("should be loaded successfully", func() {
						Expect(err).To(BeNil(), "error should be nil")
					})
					It("should be override the configuration loaded from custom configuration file", func() {
						Expect(c.Addr).To(Equal(addrCustomFlag), "c.Addr should be equal "+addrCustomFlag)
					})
					// TODO: Check Error() message
				})
			})
		})
		Context("when default system cofig file doesn't exist", func() {
			systemAddr := DEFAULT_ADDR + "0"
			systemFileContent := `addr = "` + systemAddr + `"`
			customAddr := systemAddr + "0"
			WithTempFile(systemFileContent, func(pathToSystemFile string) {
				c := New()
				c.ConfigFilePath = pathToSystemFile
				fmt.Println("======== ")
				err := c.Load([]string{"-addr", customAddr})
				It("should be loaded successfully", func() {
					Expect(err).To(BeNil(), "error should be nil")
				})
				It("should be override the default configuration loaded from default system configuration file", func() {
					Expect(c.Addr).To(Equal(customAddr), "c.Addr should be equal "+customAddr)
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
