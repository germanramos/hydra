package load_balancer_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"testing"
)

func TestBalancer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load Balancer Suite")
}
