package tests_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"testing"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests Suite")
}
