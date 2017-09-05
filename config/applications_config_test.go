package config_test

import (
	. "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"io/ioutil"
	"os"
)

var _ = Describe("ApplicationsConfig", func() {
	// TODO: Join with the config test helpers
	// HELPERS ////////////////////////////////////////////////////////////////////////
	WithTempFile := func(content string, fn func(string)) {
		f, _ := ioutil.TempFile("", "")
		f.WriteString(content)
		f.Close()
		defer os.Remove(f.Name())
		fn(f.Name())
	}
	// END OF HELPERS /////////////////////////////////////////////////////////////////
	fileContent := `[{
			"dummy1": {
				"Balancers": [
				{
					"worker": "RoundRobin",
					"simple": "OK"
				},
				{
					"worker": "SortByNumber",
					"sortAttr": "cost"
				}
				]
			}
		}, {
			"dummy2": {
				"Balancers": [
				{
					"worker": "RoundRobin",
					"simple": "OK"
				},
				{
					"worker": "SortByNumber",
					"sortAttr": "cost"
				}
				]
			}
		}]`
	Describe("Loading from JSON", func() {
		Context("When path of JSON file doesn't exist", func() {
			WithTempFile(fileContent, func(pathToFile string) {
				a := NewApplicationsConfig()
				err := a.Load(pathToFile + ".bad")
				It("should throw an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
		Context("When path of JSON file exists", func() {
			Context("When JSON is incorrect", func() {
				WithTempFile(fileContent+"???", func(pathToFile string) {
					a := NewApplicationsConfig()
					err := a.Load(pathToFile)
					It("should throw an error", func() {
						Expect(err).To(HaveOccurred())
					})
				})
			})
			Context("When JSON is correct", func() {
				WithTempFile(fileContent, func(pathToFile string) {
					a := NewApplicationsConfig()
					err := a.Load(pathToFile)
					It("should be loaded successfully", func() {
						Expect(err).To(BeNil(), "error should be nil")
						Expect(a.Apps).ToNot(BeNil())
						var apps []EtcdBaseModel
						Expect(a.Apps).To(BeAssignableToTypeOf(apps))
						apps = a.Apps
						Expect(apps).To(HaveLen(2))
						var app0 map[string]interface{}
						Expect(apps[0]).To(BeAssignableToTypeOf(app0))
						app0 = apps[0]
						Expect(app0).To(HaveKey("dummy1"))
						Expect(app0["dummy1"]).To(HaveKey("Balancers"))
						Expect(app0["dummy1"].(map[string]interface{})["Balancers"].([]interface{})[0]).To(HaveKey("worker"))
						Expect(app0["dummy1"].(map[string]interface{})["Balancers"].([]interface{})[0].(map[string]interface{})["worker"].(string)).To(Equal("RoundRobin"))
						Expect(app0["dummy1"].(map[string]interface{})["Balancers"].([]interface{})[1]).To(HaveKey("worker"))
						Expect(app0["dummy1"].(map[string]interface{})["Balancers"].([]interface{})[1].(map[string]interface{})["worker"].(string)).To(Equal("SortByNumber"))
						var app1 map[string]interface{}
						Expect(apps[1]).To(BeAssignableToTypeOf(app1))
						app1 = apps[1]
						Expect(app1).To(HaveKey("dummy2"))
						Expect(app1["dummy2"]).To(HaveKey("Balancers"))
						Expect(app1["dummy2"].(map[string]interface{})["Balancers"].([]interface{})[0]).To(HaveKey("worker"))
						Expect(app1["dummy2"].(map[string]interface{})["Balancers"].([]interface{})[0].(map[string]interface{})["worker"].(string)).To(Equal("RoundRobin"))
						Expect(app1["dummy2"].(map[string]interface{})["Balancers"].([]interface{})[1]).To(HaveKey("worker"))
						Expect(app1["dummy2"].(map[string]interface{})["Balancers"].([]interface{})[1].(map[string]interface{})["worker"].(string)).To(Equal("SortByNumber"))
					})
				})
			})
		})
	})
})
