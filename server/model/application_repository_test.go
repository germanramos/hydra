package model_test

import (
	. "github.com/innotech/hydra/server/model"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("ApplicationRepository", func() {
	Describe("Getting one application", func() {
		Context("When application ID exists", func() {
			appRepo := NewApplicationRepository()
			app := appRepo.Get("test_id")
			It("should be return an Application object", func() {
				Expect(err)
			})
		})
	})
})
