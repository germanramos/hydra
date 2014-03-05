package tests_test

import (
	. "github.com/innotech/hydra/tests/functional"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"os"

	httpUtils "github.com/innotech/hydra/vendors/github.com/coreos/etcd/tests"
)

const URL_PREFIX string = "http://127.0.0.1:8082/"

var _ = Describe("SigleNode", func() {
	Describe("Starting as standalone", func() {
		procAttr := new(os.ProcAttr)
		procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
		// args := []string{"hydra", "-name=node1", "-f", "-data-dir=/tmp/node1"}
		args := []string{"hydra", "-name=node1", "-data-dir=/tmp/node1"}

		process, err := os.StartProcess(HydraBinPath, args, procAttr)
		// process, err := os.StartProcess(HydraBinPath, args, nil)
		It("should be running successfully", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		if err != nil {
			// TODO
			GinkgoT().Fatal("start process failed:" + err.Error())
			return
		}
		defer process.Kill()

		Context("When node is running", func() {
			Context("When database is empty", func() {
				Describe("Sending a correct request to get an application", func() {
					url := URL_PREFIX + "applications/missing"
					response, err := httpUtils.Get(url, "", nil)
					It("should be got a not found response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(404))
					})
				})
				Describe("Sending a correct request to get all applications", func() {
					url := URL_PREFIX + "applications"
					response, err := httpUtils.Get(url, "", nil)
					It("should return a success response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					apps := httpUtils.ReadBodyJSON(response)
					It("should be got an empty array of json objects", func() {
						Expect(apps).To(BeEmpty())
					})
				})
				Describe("Sending a correct request to set App1 application and getting it", func() {
					app1 := map[string]interface{}{
						"App1": map[string]interface{}{
							"Cloud": "google",
						},
					}
					appJson, _ := json.Marshal(app1)
					url := URL_PREFIX + "applications"
					response, err := httpUtils.Post(url, "application/json", bytes.NewReader(appJson))
					It("should be created successfully", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					url = URL_PREFIX + "applications/App1"
					response, err = httpUtils.Get(url, "", nil)
					It("should return a success response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					apps := httpUtils.ReadBodyJSON(response)
					It("should be got an empty array of applications", func() {
						Expect(apps).To(BeEmpty())
					})
				})
			})
			Context("When App1 has been created", func() {
				Describe("Sending a correct request to override this application", func() {
					app1 := map[string]interface{}{
						"App1": map[string]interface{}{
							"Cloud": "amazon",
						},
					}
					appJson, _ := json.Marshal(app1)
					url := URL_PREFIX + "applications"
					response, err := httpUtils.Post(url, "application/json", bytes.NewReader(appJson))
					It("should be override successfully", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					url = URL_PREFIX + "applications/App1"
					response, err = httpUtils.Get(url, "", nil)
					It("should return a success response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					app := httpUtils.ReadBodyJSON(response)
					It("should be got the application overrides successfully", func() {
						Expect(app["Cloud"]).To(Equal("amazon"))
					})
				})
				Describe("Sending a correct request to set an instance in App1 application", func() {
					app1 := map[string]interface{}{
						"Instance11": map[string]interface{}{
							"Port": 8080,
						},
					}
					appJson, _ := json.Marshal(app1)
					url := URL_PREFIX + "applications"
					response, err := httpUtils.Post(url, "application/json", bytes.NewReader(appJson))
					It("should be override successfully", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					url = URL_PREFIX + "applications/App1"
					response, err = httpUtils.Get(url, "", nil)
					It("should return a success response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					app := httpUtils.ReadBodyJSON(response)
					It("should be got the application overrides successfully", func() {
						Expect(app["Port"]).To(Equal(8080))
					})
				})
			})
			Context("When App2 has not been created", func() {
				Describe("Sending a correct request to set an instance in App2 application", func() {
					app1 := map[string]interface{}{
						"Instance21": map[string]interface{}{
							"Port": 8080,
						},
					}
					appJson, _ := json.Marshal(app1)
					url := URL_PREFIX + "applications"
					response, err := httpUtils.Post(url, "application/json", bytes.NewReader(appJson))
					It("should be override successfully", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					url = URL_PREFIX + "applications/App1"
					response, err = httpUtils.Get(url, "", nil)
					It("should return a success response", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))
					})
					app := httpUtils.ReadBodyJSON(response)
					It("should be got the application overrides successfully", func() {
						Expect(app["Port"]).To(Equal(8080))
					})
				})
			})
		})
			
		})

		Context("When database is empty", func() {
			Describe("Sending a new request to set an instance in hydra application", func() {
				commit := map[string]interface{}{
					"Number": 3711,
				}
				appJson, _ := json.Marshal(commit)
				url := URL_PREFIX + "applications/hydra/instances"
				response, err := httpUtils.Post(url, "application/json", bytes.NewReader(appJson))
				It("should be created successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
			})
		})
		Context("When an application exists in database", func() {
			Describe("Sending a new request to get this application", func() {
				url := URL_PREFIX + "applications/hydra"
				response, err := httpUtils.Get(url, "application/json", nil)
				It("should be got successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
			})
		})
		Context("When an application doesn't exist in database", func() {
			Describe("Sending a correct request to get this application", func() {
				url := URL_PREFIX + "applications/missing"
				response, err := httpUtils.Get(url, "", nil)
				It("should be got successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(404))
				})
			})
		})
		Context("When some applications exist in database", func() {
			Describe("Sending a request to get these applications", func() {
				url := URL_PREFIX + "applications"
				response, err := httpUtils.Get(url, "application/json", nil)
				It("should be got successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
			})
		})
		Context("When any application don't exist in database", func() {
			Describe("Sending a correct request to get applications", func() {
				url := URL_PREFIX + "applications"
				response, err := httpUtils.Get(url, "", nil)
				It("should return a success response", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
				apps := httpUtils.ReadBodyJSON(response)
				It("should be got an empty array of json objects", func() {
					Expect(apps).To(BeEmpty())
				})
			})
		})
	})
})
