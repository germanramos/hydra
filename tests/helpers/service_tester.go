package helpers

import (
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	// "fmt"
	"log"
	"strings"
	"time"
)

type ServiceTester struct {
	baseId            string
	baseURI           string
	collection        string
	entity            string
	entityAbstractKey string
	httpUtils         *HTTPClientHelper
	serverAddr        string
}

func NewServiceTester(serverAddr, collection string, entity string, baseId string) *ServiceTester {
	s := new(ServiceTester)
	s.collection = collection
	s.baseId = baseId
	s.baseURI = "http://" + serverAddr + "/" + s.collection
	s.entity = entity
	s.entityAbstractKey = strings.Title(s.entity)
	s.httpUtils = NewHTTPClientHelper()
	s.serverAddr = serverAddr
	return s
}

func (s *ServiceTester) Pluralize(word string) string {
	return word + "s"
}

func (s *ServiceTester) DefineServiceTests(entity1 map[string]interface{}, entity2 map[string]interface{}) bool {
	return Describe(s.Pluralize(s.entityAbstractKey), func() {
		app_config := "fixtures/apps.empty.json"
		hydra_name := "hydra0"
		data_dir_path := DATA_DIR_PATH + hydra_name
		// loadBalancerAddr := "tcp://127.0.0.1:7777"
		privateAddr := "127.0.0.1:7771"
		// publicAddr := "127.0.0.1:7772"
		// pingInstancesAddr := "http://" + privateAddr + "/apps/Ping/Instances"
		args := []string{"-name=" + hydra_name, "-private-addr=" + privateAddr, "-data-dir=" + data_dir_path, "-apps-file=" + app_config}
		process := RunHydraInStandaloneAndReturnProcess(args)
		defer KillHydraProcess(process)
		time.Sleep(time.Second)

		Context("When database is empty", func() {
			Describe("Sending a correct request to get a missing "+s.entity, func() {
				response, getError := s.httpUtils.Get(s.baseURI + "/" + s.entityAbstractKey + "1")
				It("should be got a not found response", func() {
					Expect(getError).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(404))
				})
			})
			Describe("Sending a correct request to get all "+s.collection, func() {
				response, getAllError := s.httpUtils.Get(s.baseURI)
				It("should be got a not found response", func() {
					Expect(getAllError).To(BeNil())
					Expect(response.StatusCode).To(Equal(404))
				})
			})
			appJson, _ := json.Marshal(entity1)
			Describe("Sending a bad request to set "+s.entityAbstractKey+"1 application", func() {
				badAppJson := appJson[5:]
				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(badAppJson))
				It("should return a bad request status code", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(400))
				})
			})
			Describe("Setting "+s.entityAbstractKey+"1 application and getting it", func() {
				response1, err1 := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
				It("should be created successfully", func() {
					Expect(err1).NotTo(HaveOccurred())
					Expect(response1.StatusCode).To(Equal(200))
				})
				log.Print("********************" + s.baseURI + "/" + s.baseId + "1")
				response2, err2 := s.httpUtils.Get(s.baseURI + "/" + s.baseId + "1")
				log.Print(response2)
				It("should return a success response", func() {
					Expect(err2).NotTo(HaveOccurred())
					Expect(response2.StatusCode).To(Equal(200))
				})
				entity, err3 := s.httpUtils.ReadBodyJsonObject(response2)
				It("should be got the correct application", func() {
					Expect(err3).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
					Expect(entity).NotTo(BeNil())
					Expect(entity).To(Equal(entity1))
				})
			})
		})
		Context("When App1 have been created", func() {
			appJson, _ := json.Marshal(entity1)
			Describe("Overriding "+s.entityAbstractKey+"1 application", func() {
				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
				It("should be created successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
				response, err = s.httpUtils.Get(s.baseURI + "/" + s.baseId + "1")
				It("should return a success response", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
				entity, err := s.httpUtils.ReadBodyJsonObject(response)
				It("should be got the correct application", func() {
					Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
					Expect(entity).NotTo(BeEmpty())
					Expect(entity).To(Equal(entity1))
				})
			})
			appJson, _ = json.Marshal(entity2)
			Describe("Setting "+s.entityAbstractKey+"2 application and getting all applications ("+s.entityAbstractKey+"1 and "+s.entityAbstractKey+"2)", func() {
				response, err := s.httpUtils.Post(s.baseURI, "application/json", bytes.NewReader(appJson))
				It("should be created successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
				response, err = s.httpUtils.Get(s.baseURI)
				It("should return a success response", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))
				})
				entities, err := s.httpUtils.ReadBodyJsonArray(response)
				It("should be got the correct application", func() {
					Expect(err).NotTo(HaveOccurred(), "HTTP body JSON should be a valid json")
					Expect(entities).NotTo(BeEmpty())
					Expect(entities).To(HaveLen(2))
					Expect(entities[0]).To(Equal(entity1))
					Expect(entities[1]).To(Equal(entity2))
				})
			})
		})
	})
}
