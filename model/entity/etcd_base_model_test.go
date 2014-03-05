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
			GinkgoT().Fatal("JSON blob wrongly defined")
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

	Describe("Intantiating new etcd model from etcd store event", func() {
		var event = &store.Event{
			Action: "get",
			Node: &store.NodeExtern{
				Key: "/",
				Dir: true,
				Nodes: store.NodeExterns{
					&store.NodeExtern{
						Key: "key_1",
						Dir: true,
						Nodes: store.NodeExterns{
							&store.NodeExtern{
								Key:           "key_1_1",
								Value:         "24.00",
								ModifiedIndex: 8,
								CreatedIndex:  8,
							},
							&store.NodeExtern{
								Key:           "key_1_2",
								Value:         "true",
								ModifiedIndex: 8,
								CreatedIndex:  8,
							},
						},
					},
					&store.NodeExtern{
						Key: "key_2",
						Dir: true,
						Nodes: store.NodeExterns{
							&store.NodeExtern{
								Key:           "key_2_1",
								Value:         "Hello Hydra",
								ModifiedIndex: 12,
								CreatedIndex:  12,
							},
							&store.NodeExtern{
								Key:           "key_2_2",
								Value:         "",
								ModifiedIndex: 12,
								CreatedIndex:  12,
							},
						},
					},
				},
			},
		}

		m1, err := entity.NewFromEvent(event)
		m := map[string]interface{}(*m1)
		It("should instantiate a new EtcdBaseModel successfully", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(m["/"].(map[string]interface{})["key_1"].(map[string]interface{})["key_1_1"].(string)).To(Equal("24.00"))
			Expect(m["/"].(map[string]interface{})["key_1"].(map[string]interface{})["key_1_2"]).To(Equal("true"))
			Expect(m["/"].(map[string]interface{})["key_2"].(map[string]interface{})["key_2_1"]).To(Equal("Hello Hydra"))
			Expect(m["/"].(map[string]interface{})["key_2"].(map[string]interface{})["key_2_2"]).To(Equal(""))
		})
	})

	// Describe("Checking if i field exits", func() {
	// 	Context("When correct structure instance contains the i field", func() {
	// 		type myStruct struct {
	// 			i int
	// 			j int
	// 			s string
	// 		}
	// 		var s myStruct
	// 		exists, err := entity.CheckIfStructFieldNameExists(s, "i")
	// 		It("should exist", func() {
	// 			Expect(err).NotTo(HaveOccurred())
	// 			Expect(exists).To(BeTrue())
	// 		})
	// 	})

	// 	Context("When correct structure instance doesn't contain the i field", func() {
	// 		type myStruct struct {
	// 			j int
	// 			s string
	// 		}
	// 		var s myStruct
	// 		exists, err := entity.CheckIfStructFieldNameExists(s, "i")
	// 		It("should exist", func() {
	// 			Expect(err).NotTo(HaveOccurred())
	// 			Expect(exists).To(BeFalse())
	// 		})
	// 	})

	// 	Context("When no structure instance contains the i field", func() {
	// 		var s int
	// 		exists, err := entity.CheckIfStructFieldNameExists(s, "i")
	// 		It("should exist", func() {
	// 			Expect(err).To(HaveOccurred())
	// 		})
	// 	})

	// })
})
