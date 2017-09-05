package helpers

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

var FIXTURES_PATH string

var HydraBinPath string

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	FIXTURES_PATH = pwd

	// Initialize the 'hydra' binary path or default it to the hydra diretory.
	HydraBinPath = os.Getenv("HYDRA_BIN_PATH")
	if HydraBinPath == "" {
		HydraBinPath = filepath.Join(build.Default.GOPATH, "src", "github.com", "innotech", "hydra", "hydra")
	}
}
