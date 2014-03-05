package application_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"testing"

	. "github.com/innotech/hydra/tests/helpers"
)

func TestApplication(t *testing.T) {
	process := RunHydraInStandaloneAndReturnProcess()
	// defer process.kill()
	defer KillHydraProcess(process)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Suite")
}
