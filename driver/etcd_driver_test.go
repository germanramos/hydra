package driver_test

import (
	. "github.com/innotech/hydra/driver"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"github.com/innotech/hydra/server/model"
	etcdMocks "github.com/innotech/hydra/vendors/github.com/coreos/etcd/tests/mock"
)

var _ = Describe("EtcdDriver", func() {
	Describe("Creating new entity", func() {
		Context("When entity is valid", func() {
			storeMock := etcdMocks.NewStore()
			driver := NewEtcdDriver(storeMock)
			err := driver.Create("")
			It("should be persisted successfully", func() {
				// Expect(err)
			})
		})
	})
})
