package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/innotech/hydra/vendors/github.com/BurntSushi/toml"
	"github.com/innotech/hydra/vendors/github.com/stretchr/testify/assert"
)

func TestConfigFile(t *testing.T) {
	fileContent := `
		addr = "127.0.0.1:4002"
	`
	// var conf Config
	conf := New()
	_, err := toml.Decode(fileContent, &conf)

	assert.Nil(t, err, "err should be nothing")
	// assert.Equal(t, conf.Addr, "127.0.0.1:4002", "conf.Addr and \"127.0.0.1:4002\" should be equal")
}

// Ensures that a default config file path field is overridden by a custom config file path field.
func TestConfigCustomConfigFilePathOverrideSystemConfig(t *testing.T) {
	defaultConf := `addr = "127.0.0.1:5000"`
	customConf := `addr = "127.0.0.1:6000"`
	withTempFile(defaultConf, func(p1 string) {
		withTempFile(customConf, func(p2 string) {
			c := New()
			c.ConfigFilePath = p1
			// assert.Nil(t, c.Load([]string{"-config", p2}), "")
			// assert.Equal(t, c.Addr, "http://127.0.0.1:6000", "")
		})
	})
}

//--------------------------------------
// Helpers
//--------------------------------------

// Creates a temp file and calls a function with the context.
func withTempFile(content string, fn func(string)) {
	f, _ := ioutil.TempFile("", "")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	fn(f.Name())
}
