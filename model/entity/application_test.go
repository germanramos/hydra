package entity_test

import (
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	Describe("Creating new application", func() {
		Context("When arguments are correct", func() {
			appData := map[string]interface{}{
				"Balancers": map[string]interface{}{
					"Ping": map[string]interface{}{
						"Mode": "Direct",
					},
					"Pong": map[string]interface{}{
						"Incr": "44.21",
					},
				},
				"Instances": map[string]interface{}{
					"Instance-1": map[string]interface{}{
						"cpu": "33.67",
					},
					"Instance-2": map[string]interface{}{
						"mem": "55.78",
					},
				},
			}
			app, err := NewApplication("app-1", appData)
			It("should not throw an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should contain the correct Id", func() {
				Expect(app.Id).To(Equal("app-1"))
			})
			It("should contain the correct Balancers", func() {
				Expect(app.Balancers).To(HaveLen(2))
				Expect(app.Balancers[0].Id).To(Equal("Ping"))
				Expect(app.Balancers[0].Args["Mode"]).To(Equal("Direct"))
				Expect(app.Balancers[1].Id).To(Equal("Pong"))
				Expect(app.Balancers[1].Args["Incr"]).To(Equal("44.21"))
			})
			It("should contain the correct Instances", func() {
				Expect(app.Instances).To(HaveLen(2))
				Expect(app.Instances[0].Id).To(Equal("Instance-1"))
				Expect(app.Instances[0].Info["cpu"]).To(Equal("33.67"))
				Expect(app.Instances[1].Id).To(Equal("Instance-2"))
				Expect(app.Instances[1].Info["mem"]).To(Equal("55.78"))
			})
		})
	})
})
