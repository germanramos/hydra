package balancer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBalancer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Balancer Suite")
}
