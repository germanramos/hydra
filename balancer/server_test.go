package balancer_test

import (
	. "github.com/innotech/hydra/balancer"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Describe("Registering application plumber", func() {
		b := NewServer()
		appId, appAttrs := 
		b.RegisterPlumber(app)
		It("should be registered succesfully", func() {
			Expect(b, ...)
			
		})
	})
})
