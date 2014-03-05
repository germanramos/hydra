package application_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	"encoding/json"

	. "github.com/innotech/hydra/tests/helpers"
	// . "github.com/innotech/hydra/tests/helpers"
)

const BASE_URI string = HYDRA_URI + "/applications"

var _ = Describe("Applications", func() {
	httpUtils := NewHTTPClientHelper()
	Context("When database is empty", func() {
		Describe("Sending a correct request to get a missing application", func() {
			response, err := httpUtils.Get(BASE_URI + "/App1")
			It("should be got a not found response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(404))
			})
		})
		Describe("Sending a correct request to get all applications", func() {
			response, err := httpUtils.Get(BASE_URI)
			It("should return a success response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			apps, err := httpUtils.ReadBodyJSON(response)
			It("should be got an empty array of json objects", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(apps).To(BeEmpty())
			})
		})
		app1 := map[string]interface{}{
			"App1": map[string]interface{}{
				"Cloud": "google",
			},
		}
		appJson, _ := json.Marshal(app1)
		Describe("Sending a bad request to set App1 application", func() {
			badAppJson := appJson[5:]
			response, err := httpUtils.Post(BASE_URI, "application/json", bytes.NewReader(badAppJson))
			It("should return a bad request status code", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(400))
			})
		})
		Describe("Setting App1 application and getting it", func() {
			response, err := httpUtils.Post(BASE_URI, "application/json", bytes.NewReader(appJson))
			It("should be created successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			response, err = httpUtils.Get(BASE_URI + "/App1")
			It("should return a success response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			app, err := httpUtils.ReadBodyJSON(response)
			It("should be got the correct application", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(app).NotTo(BeEmpty())
				Expect(app["Cloud"]).To(Equal("google"))
			})
		})
	})
	Context("When App1 have been created", func() {
		app1 := map[string]interface{}{
			"App1": map[string]interface{}{
				"Cloud": "amazon",
			},
		}
		appJson, _ := json.Marshal(app1)
		Describe("Overriding App1 application", func() {
			response, err := httpUtils.Post(BASE_URI, "application/json", bytes.NewReader(appJson))
			It("should be created successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			response, err = httpUtils.Get(BASE_URI + "/App1")
			It("should return a success response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			app, err := httpUtils.ReadBodyJSON(response)
			It("should be got the correct application", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(app).NotTo(BeEmpty())
				Expect(app["Cloud"]).To(Equal("amazon"))
			})
		})
		app2 := map[string]interface{}{
			"App2": map[string]interface{}{
				"Cloud": "azure",
			},
		}
		appJson, _ = json.Marshal(app2)
		Describe("Setting App2 application and getting all applications (App1 and App2)", func() {
			response, err := httpUtils.Post(BASE_URI, "application/json", bytes.NewReader(appJson))
			It("should be created successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			response, err = httpUtils.Get(BASE_URI)
			It("should return a success response", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
			})
			apps, err := httpUtils.ReadBodyJSON(response)
			It("should be got the correct application", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(apps).NotTo(BeEmpty())
				Expect(apps).NotTo(HaveLen(2))
				Expect(apps["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("amazon"))
				Expect(apps["App2"].(map[string]interface{})["Cloud"].(string)).To(Equal("azure"))
			})
		})
	})
})
