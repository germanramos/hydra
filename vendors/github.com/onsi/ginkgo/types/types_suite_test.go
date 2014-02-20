package types_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgo Types Suite")
}
