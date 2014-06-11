package controller_test

import (
	. "github.com/innotech/hydra/server/controller"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("BasicController", func() {
	const (
		VAR_0 string = "appId"
		VAR_1 string = "instanceId"
	)
	Describe("Making new instance", func() {
		Context("When base path contains incomplete variable placeholder", func() {
			Context("At the beginning of the placeholder", func() {
				basePath := "/applications/" + VAR_0 + "}/instances"
				b, err := NewBasicController(basePath)
				It("should be extracted the path variables", func() {
					Expect(err).To(HaveOccurred())
					Expect(b).To(BeNil())
				})
			})
			Context("At the ending of the placeholder", func() {
				basePath := "/applications/{" + VAR_0 + "/instances"
				b, err := NewBasicController(basePath)
				It("should be extracted the path variables", func() {
					Expect(err).To(HaveOccurred())
					Expect(b).To(BeNil())
				})
			})
		})
		Context("When base path contains no variables", func() {
			basePath := "/applications"
			b, err := NewBasicController(basePath)
			It("should be extracted the path variables", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(b.PathVariables).To(HaveLen(0))
			})
		})
		Context("When base path contains variables", func() {
			basePath := "/applications/{" + VAR_0 + "}/instances/{" + VAR_1 + "}/stats"
			b, err := NewBasicController(basePath)
			It("should be extracted the path variables", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(b.PathVariables).ToNot(BeNil())
				Expect(b.PathVariables).To(HaveLen(2))
				Expect(b.PathVariables[0]).To(Equal(VAR_0))
				Expect(b.PathVariables[1]).To(Equal(VAR_1))
			})
		})
	})
	Describe("Getting repository from path variables of request", func() {
		basePath := "/applications/{" + VAR_0 + "}/instances"
		b, _ := NewBasicController(basePath)
		vars := make(map[string]string)
		vars[VAR_0] = "app_1"
		vars["id"] = "instance_1"
		repo := b.GetConfiguredRepository(vars)
		It("should be returned the repository properly configured", func() {
			Expect(repo.GetCollection()).To(Equal("/applications/app_1/instances"))
		})
	})
})
