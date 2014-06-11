package subpackage

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("nested sub packages", func() {
		GinkgoT().Fail(true)
	})
})
