package public_api_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"testing"
)

func TestPublic_api(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public_api Suite")
}
