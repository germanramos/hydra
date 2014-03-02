package entity_test

import (
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"encoding/json"
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
})
