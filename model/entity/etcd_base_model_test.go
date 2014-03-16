package entity_test

import (
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"encoding/json"

	"github.com/innotech/hydra/model/entity"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
)

var _ = Describe("EtcdBaseModel", func() {
	Describe("Exploding", func() {
		Context("When the entity is empty", func() {
			e := new(EtcdBaseModel)
			id, data := e.Explode()
			It("should be returned an empty id", func() {
				Expect(id).To(BeEmpty())
			})
			It("should not return data", func() {
				Expect(data).To(BeNil())
			})
		})
		Context("When the entity has more than one super key", func() {
			var jsonBlob = []byte(`{
				"Number": 12.657,
				"Bool": true
			}`)
			var e EtcdBaseModel
			err := json.Unmarshal(jsonBlob, &e)
			if err != nil {
				Fail("JSON blob wrongly defined")
			}
			id, data := e.Explode()
			It("should be returned an empty id", func() {
				Expect(id).To(BeEmpty())
			})
			It("should not return data", func() {
				Expect(data).To(BeNil())
			})
		})
		Context("When the entity has only one super key as entity id", func() {
			var jsonBlob = []byte(`{
				"entityID": {
					"Number": 12.65,
					"Bool": true
				}
			}`)
			var e EtcdBaseModel
			err := json.Unmarshal(jsonBlob, &e)
			if err != nil {
				Fail("JSON blob wrongly defined")
			}
			id, data := e.Explode()

			It("should be returned an empty id", func() {
				Expect(id).To(Equal("entityID"))
			})
			It("should not return data", func() {
				Expect(data).ToNot(BeNil())
				Expect(data).To(HaveLen(2))
				Expect(data).To(HaveKey("Number"))
				Expect(data["Number"]).To(Equal(12.65))
				Expect(data).To(HaveKey("Bool"))
				Expect(data["Bool"]).To(Equal(true))
			})
		})
	})
	Describe("Exporting etcd operations", func() {
		var jsonBlob = []byte(`{
			"Null": null,
			"Number": 12.657,
			"Bool": true,
			"String": "Welcome to Hydra",
			"Array": ["element0", "element1"],
			"Object": {
				"Number1": 55,
				"String1": "Choose your strategy"
			},
			"ArrayOfObjects": [{"ObjectKey0": false}, {"ObjectKey0": "And start your broker"}]
		}`)
		var e EtcdBaseModel
		err := json.Unmarshal(jsonBlob, &e)
		if err != nil {
			Fail("JSON blob wrongly defined")
			// GinkgoT().Fatal("JSON blob wrongly defined")
		}
		etcdOps, err := e.ExportEtcdOperations()
		It("should be exported successfully", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(etcdOps["/Null"]).To(Equal(""))
			Expect(etcdOps["/Number"]).To(Equal("12.66"))
			Expect(etcdOps["/Bool"]).To(Equal("true"))
			Expect(etcdOps["/String"]).To(Equal("Welcome to Hydra"))
			Expect(etcdOps["/Array/0"]).To(Equal("element0"))
			Expect(etcdOps["/Array/1"]).To(Equal("element1"))
			Expect(etcdOps["/Object/Number1"]).To(Equal("55.00"))
			Expect(etcdOps["/Object/String1"]).To(Equal("Choose your strategy"))
			Expect(etcdOps["/ArrayOfObjects/0/ObjectKey0"]).To(Equal("false"))
			Expect(etcdOps["/ArrayOfObjects/1/ObjectKey0"]).To(Equal("And start your broker"))
		})
	})

	Describe("Extracting json key from etcd key", func() {
		Context("When json key is empty", func() {
			jsonKey, err := ExtractJsonKeyFromEtcdKey("")
			It("should be thrown an error", func() {
				Expect(err).To(HaveOccurred())
				// TODO: check error type or message
				Expect(jsonKey).To(BeEmpty())
			})
		})
		Context("When json key is incorrect", func() {
			jsonKey, err := ExtractJsonKeyFromEtcdKey("/db/applications/")
			It("should be thrown an error", func() {
				Expect(err).To(HaveOccurred())
				// TODO: check error type or message
				Expect(jsonKey).To(BeEmpty())
			})
		})
		Context("When json key is correct", func() {
			jsonKey, err := ExtractJsonKeyFromEtcdKey("/db/applications/App1/instances/Instance1")
			It("should be extracted successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(jsonKey).To(Equal("Instance1"))
			})
		})
	})

	Context("When store is not empty", func() {
		event := &store.Event{
			Action: "get",
			Node: &store.NodeExtern{
				Key: "/App1",
				Dir: true,
				Nodes: store.NodeExterns{
					&store.NodeExtern{
						Key: "/App1/key_1",
						Dir: true,
						Nodes: store.NodeExterns{
							&store.NodeExtern{
								Key:           "/App1/key_1/key_1_1",
								Value:         "24.00",
								ModifiedIndex: 8,
								CreatedIndex:  8,
							},
							&store.NodeExtern{
								Key:           "/App1/key_1/key_1_2",
								Value:         "true",
								ModifiedIndex: 8,
								CreatedIndex:  8,
							},
						},
					},
					&store.NodeExtern{
						Key: "/App1/key_2",
						Dir: true,
						Nodes: store.NodeExterns{
							&store.NodeExtern{
								Key:           "/App1/key_2/key_2_1",
								Value:         "Hello Hydra",
								ModifiedIndex: 12,
								CreatedIndex:  12,
							},
							&store.NodeExtern{
								Key:           "/App1/key_2/key_2_2",
								Value:         "",
								ModifiedIndex: 12,
								CreatedIndex:  12,
							},
						},
					},
				},
			},
		}
		Describe("Instantiating new etcd model from etcd store event", func() {
			m1, err := entity.NewModelFromEvent(event)
			m := map[string]interface{}(*m1)
			It("should instantiate a new EtcdBaseModel successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(m["App1"].(map[string]interface{})["key_1"].(map[string]interface{})["key_1_1"].(string)).To(Equal("24.00"))
				Expect(m["App1"].(map[string]interface{})["key_1"].(map[string]interface{})["key_1_2"]).To(Equal("true"))
				Expect(m["App1"].(map[string]interface{})["key_2"].(map[string]interface{})["key_2_1"]).To(Equal("Hello Hydra"))
				Expect(m["App1"].(map[string]interface{})["key_2"].(map[string]interface{})["key_2_2"]).To(Equal(""))
			})
		})
		Describe("Instantiating new etcd models (array of models) from etcd store event", func() {
			m1, err := entity.NewModelsFromEvent(event)
			m := []EtcdBaseModel(*m1)
			It("should instantiate a new EtcdBaseModel successfully", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(m[0]["key_1"].(map[string]interface{})["key_1_1"].(string)).To(Equal("24.00"))
				Expect(m[0]["key_1"].(map[string]interface{})["key_1_2"]).To(Equal("true"))
				Expect(m[1]["key_2"].(map[string]interface{})["key_2_1"]).To(Equal("Hello Hydra"))
				Expect(m[1]["key_2"].(map[string]interface{})["key_2_2"]).To(Equal(""))
			})
		})
	})
})
