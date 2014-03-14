package balancer_test

import (
	. "github.com/innotech/hydra/balancer"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("Balancer", func() {
	Describe("Receiving request", func() {
		b := NewBalancer()
		app1 := map[string]interface{}{
			"App-1": map[string]interface{}{
				"Cloud": "amazon",
				"Instances": map[string]map[string]interface{}{
					"Instance-1": {
						"cpu": 80.00,
						"mem": 40.00,
					},
					"Instance-2": {
						"cpu": 30.00,
						"mem": 70.00,
					},
				},
			},
		}
		b.Balance(app1)
		// TODO: Expects
	})
})
