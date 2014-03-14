package balancer_test

import (
	. "github.com/innotech/hydra/balancer"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"github.com/innotech/hydra/config"
)

var _ = Describe("Plumber", func() {
	Describe("Cofiguring", func() {
		Context("When correct information has been sent", func() {
			// Defining Config Mock
			c := config.New()
			c.Balancers = make(map[string][]map[string]interface{})
			c.Balancers["app-1"] = []map[string]map[string]interface{}{
				map[string]interface{}{
					"id": "cloud-map",
				},
				map[string]interface{}{
					"id": "cpu-load",
				},
			}
			c.Balancers["app-2"] = []map[string]map[string]interface{}{
				map[string]interface{}{
					"id": "cloud-map",
				},
				map[string]interface{}{
					"id": "mem-load",
				},
			}
			// Loading Plumber
			p := NewPlumber()
			p.Configure(c.getBalancers())
			// TODO: Expects
		})
		Context("When incorrect information has been sent", func() {
			c := config.New()
		})
	})
})
