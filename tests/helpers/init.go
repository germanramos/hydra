package helpers

import (
	"go/build"
	"os"
	"path/filepath"
)

var HydraBinPath string

func init() {
	// Initialize the 'hydra' binary path or default it to the hydra diretory.
	HydraBinPath = os.Getenv("HYDRA_BIN_PATH")
	if HydraBinPath == "" {
		HydraBinPath = filepath.Join(build.Default.GOPATH, "src", "github.com", "innotech", "hydra", "hydra")
	}
}
