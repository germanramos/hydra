package config

import (
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
	assert.Equal(t, conf.Addr, "127.0.0.1:4002", "conf.Addr and \"127.0.0.1:4002\" should be equal")
}
