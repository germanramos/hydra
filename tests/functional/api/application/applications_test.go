package application_test

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"time"

	. "github.com/innotech/hydra/tests/helpers"
)

const BASE_URI string = "http://" + PRIVATE_HYDRA_URI + "/applications"

var _ = Describe("Applications", func() {
	process := RunHydraInStandaloneAndReturnProcess()
	defer KillHydraProcess(process)
	time.Sleep(time.Second)

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
			response, err1 := httpUtils.Get(BASE_URI)
			It("should return a success response", func() {
				Expect(err1).To(BeNil())
				Expect(response.StatusCode).To(Equal(404))
			})
			// apps, err2 := httpUtils.ReadBodyJSON(response)
			// It("should be got an empty array of json objects", func() {
			// 	Expect(err2).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
			// 	Expect(apps).To(BeEmpty())
			// })
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
			response1, err1 := httpUtils.Post(BASE_URI, "application/json", bytes.NewReader(appJson))
			It("should be created successfully", func() {
				Expect(err1).NotTo(HaveOccurred())
				Expect(response1.StatusCode).To(Equal(200))
			})
			response2, err2 := httpUtils.Get(BASE_URI + "/App1")
			It("should return a success response", func() {
				Expect(err2).NotTo(HaveOccurred())
				Expect(response2.StatusCode).To(Equal(200))
			})
			app, err3 := httpUtils.ReadBodyJsonObject(response2)
			It("should be got the correct application", func() {
				Expect(err3).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(app).NotTo(BeNil())
				Expect(app["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("google"))
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
			app, err := httpUtils.ReadBodyJsonObject(response)
			It("should be got the correct application", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(app).NotTo(BeEmpty())
				Expect(app["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("amazon"))
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
			apps, err := httpUtils.ReadBodyJsonArray(response)
			It("should be got the correct application", func() {
				Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
				Expect(apps).NotTo(BeEmpty())
				Expect(apps).To(HaveLen(2))
				Expect(apps[0]["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("amazon"))
				Expect(apps[1]["App2"].(map[string]interface{})["Cloud"].(string)).To(Equal("azure"))
			})
		})
	})
})
