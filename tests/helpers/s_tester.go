package helpers

// import (
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	"bytes"
// 	"encoding/json"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// type ServiceTester struct {
// 	baseURI       string
// 	collection    string
// 	entity        string
// 	entityCounter int
// 	httpUtils     *HTTPClientHelper
// 	serverAddr    string
// }

// func NewServiceTester(serverAddr, entity string /*, parentEntities []string*/) *ServiceTester {
// 	s := new(ServiceTester)
// 	s.baseURI = "http://" + serverAddr + "/" + Pluralize(entity)
// 	// s.baseURI = "http://" + serverAddr + "/"
// 	s.entity = entity
// 	s.entityCounter = 0
// 	s.httpUtils = NewHTTPClientHelper()
// 	s.serverAddr = serverAddr
// 	return s
// }

// func Pluralize(word string) string {
// 	return word + "s"
// }

// // func fake(s *ServiceTester) {
// // 	Describe("Sending a correct request to get a missing "+s.entity, func() {
// // 		response, err := s.httpUtils.Get(s.baseURI + "/" + strings.Title(s.entity) + "1")
// // 		It("should be got a not found response", func() {
// // 			Expect(err).NotTo(HaveOccurred())
// // 			Expect(response.StatusCode).To(Equal(404))
// // 		})
// // 	})
// // }

// func (s *ServiceTester) generateId() string {
// 	s.entityCounter++
// 	return strings.Title(s.entity) + strconv.Itoa(s.entityCounter)
// }

// func (s *ServiceTester) DefineServiceTests() bool {
// 	return Describe(Pluralize(strings.Title(s.entity)), func() {
// 		process := RunHydraInStandaloneAndReturnProcess(s.serverAddr)
// 		defer KillHydraProcess(process)
// 		time.Sleep(time.Second)

// 		Context("When database is empty", func() {
// 			Describe("Sending a correct request to get a missing "+s.entity, func() {
// 				response, getError := s.httpUtils.Get(s.baseURI + "/" + s.generateId())
// 				It("should be got a not found response", func() {
// 					Expect(getError).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(404))
// 				})
// 			})
// 			// fake(s)
// 			Describe("Sending a correct request to get all "+s.collection, func() {
// 				response, getAllError := s.httpUtils.Get(s.baseURI)
// 				It("should be got a not found response", func() {
// 					Expect(getAllError).To(BeNil())
// 					Expect(response.StatusCode).To(Equal(404))
// 				})
// 			})
// 			app1 := map[string]interface{}{
// 				"App1": map[string]interface{}{
// 					"Cloud": "google",
// 				},
// 			}
// 			appJson, _ := json.Marshal(app1)
// 			Describe("Sending a bad request to set App1 application", func() {
// 				badAppJson := appJson[5:]
// 				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(badAppJson))
// 				It("should return a bad request status code", func() {
// 					Expect(err).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(400))
// 				})
// 			})
// 			Describe("Setting App1 application and getting it", func() {
// 				response1, err1 := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
// 				It("should be created successfully", func() {
// 					Expect(err1).NotTo(HaveOccurred())
// 					Expect(response1.StatusCode).To(Equal(200))
// 				})
// 				response2, err2 := s.httpUtils.Get(s.baseURI + "/App1")
// 				It("should return a success response", func() {
// 					Expect(err2).NotTo(HaveOccurred())
// 					Expect(response2.StatusCode).To(Equal(200))
// 				})
// 				app, err3 := s.httpUtils.ReadBodyJsonObject(response2)
// 				It("should be got the correct application", func() {
// 					Expect(err3).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
// 					Expect(app).NotTo(BeNil())
// 					Expect(app["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("google"))
// 				})
// 			})
// 		})
// 		Context("When App1 have been created", func() {
// 			app1 := map[string]interface{}{
// 				"App1": map[string]interface{}{
// 					"Cloud": "amazon",
// 				},
// 			}
// 			appJson, _ := json.Marshal(app1)
// 			Describe("Overriding App1 application", func() {
// 				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
// 				It("should be created successfully", func() {
// 					Expect(err).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(200))
// 				})
// 				response, err = s.httpUtils.Get(s.baseURI + "/App1")
// 				It("should return a success response", func() {
// 					Expect(err).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(200))
// 				})
// 				app, err := s.httpUtils.ReadBodyJsonObject(response)
// 				It("should be got the correct application", func() {
// 					Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
// 					Expect(app).NotTo(BeEmpty())
// 					Expect(app["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("amazon"))
// 				})
// 			})
// 			app2 := map[string]interface{}{
// 				"App2": map[string]interface{}{
// 					"Cloud": "azure",
// 				},
// 			}
// 			appJson, _ = json.Marshal(app2)
// 			Describe("Setting App2 application and getting all applications (App1 and App2)", func() {
// 				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
// 				It("should be created successfully", func() {
// 					Expect(err).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(200))
// 				})
// 				response, err = s.httpUtils.Get(s.baseURI)
// 				It("should return a success response", func() {
// 					Expect(err).NotTo(HaveOccurred())
// 					Expect(response.StatusCode).To(Equal(200))
// 				})
// 				apps, err := s.httpUtils.ReadBodyJsonArray(response)
// 				It("should be got the correct application", func() {
// 					Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
// 					Expect(apps).NotTo(BeEmpty())
// 					Expect(apps).To(HaveLen(2))
// 					Expect(apps[0]["App1"].(map[string]interface{})["Cloud"].(string)).To(Equal("amazon"))
// 					Expect(apps[1]["App2"].(map[string]interface{})["Cloud"].(string)).To(Equal("azure"))
// 				})
// 			})
// 		})
// 	})
// }
