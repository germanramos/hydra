package config_test

import (
	. "github.com/innotech/hydra/config"
	. "github.com/innotech/hydra/model/entity"
	"github.com/innotech/hydra/model/repository/mock_repository"
	"github.com/innotech/hydra/vendors/code.google.com/p/gomock/gomock"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	"github.com/innotech/hydra/vendors/github.com/onsi/ginkgo/thirdparty/gomocktestreporter"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"fmt"
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
 				"balancers": {
 					"cloud-map": {
 					},
					"cpu-load": {
					}
				}
			}
		}, {
			"dummy2": {
 				"balancers": {
 					"cloud-map": {
					},
					"mem-load": {
					}
				}
			}
		}]`
	Describe("Loading from JSON", func() {
		// fileContent := `[{
		// 	"dummy1": {
		// 			"balancers": {
		// 				"cloud-map": {
		// 				},
		// 			"cpu-load": {
		// 			}
		// 		}
		// 	}
		// }, {
		// 	"dummy2": {
		// 			"balancers": {
		// 				"cloud-map": {
		// 			},
		// 			"mem-load": {
		// 			}
		// 		}
		// 	}
		// }]`
		Context("When path of JSON file doesn't exist", func() {
			WithTempFile(fileContent, func(pathToFile string) {
				a := NewApplicationsConfig()
				err := a.Load(pathToFile + ".bad")
				It("should throw an error", func() {
					Expect(err).To(HaveOccurred())
					// TODO: check kind of error
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
						// TODO: check kind of error
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
						Expect(app0["dummy1"]).To(HaveKey("balancers"))
						Expect(app0["dummy1"].(map[string]interface{})["balancers"]).To(HaveKey("cloud-map"))
						Expect(app0["dummy1"].(map[string]interface{})["balancers"]).To(HaveKey("cpu-load"))
						var app1 map[string]interface{}
						Expect(apps[1]).To(BeAssignableToTypeOf(app1))
						app1 = apps[1]
						Expect(app1).To(HaveKey("dummy2"))
						Expect(app1["dummy2"]).To(HaveKey("balancers"))
						Expect(app1["dummy2"].(map[string]interface{})["balancers"]).To(HaveKey("cloud-map"))
						Expect(app1["dummy2"].(map[string]interface{})["balancers"]).To(HaveKey("mem-load"))
					})
				})
			})
		})
	})
	Describe("Saving loaded applications", func() {
		var (
			mockCtrl  *gomock.Controller
			mockRepo  *mock_repository.MockEtcdAccessLayer
			appConfig *ApplicationsConfig
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(gomocktestreporter.New())
			// mockThing = mockthing.NewMockThing(mockCtrl)
			mockRepo = mock_repository.NewMockEtcdAccessLayer(mockCtrl)
			appConfig = NewApplicationsConfig()
			// repo := appConfig.Repo
			// appConfig.Repo = mockRepo
			// appConfig.Repo = repo
		})

		AfterEach(func() {
			mockCtrl.Finish()
		})

		It("should persist applications successfully", func() {
			// var previusCall *gomock.Call
			// hasPreviusCall := false
			// for _, app := range appConfig.Apps {
			// 	if !hasPreviusCall {
			// 		previusCall = mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
			// 		hasPreviusCall = true
			// 	} else {
			// 		previusCall = mockRepo.EXPECT().Set(app).After(previusCall)
			// 	}
			// }
			appConfig = NewApplicationsConfig()
			appConfig.Repo = mockRepo
			Expect(appConfig.Repo).To(Equal(mockRepo))

			// mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
			err := appConfig.Persists()
			// mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		WithTempFile(fileContent, func(pathToFile string) {
			mockCtrl = gomock.NewController(gomocktestreporter.New())
			defer mockCtrl.Finish()
			// mockThing = mockthing.NewMockThing(mockCtrl)
			mockRepo = mock_repository.NewMockEtcdAccessLayer(mockCtrl)
			// mockRepo.setCollection("applications")
			appConfig = NewApplicationsConfig()
			// appConfig.Repo = mockRepo
			// repo := appConfig.Repo
			// fmt.Println(pathToFile)
			// a := NewApplicationsConfig()
			// err := a.Load(pathToFile)
			err := appConfig.Load(pathToFile)
			fmt.Println("WWW")
			if err != nil {
				fmt.Println(err)
			}
			// var err error = nil
			It("should persist applications successfully", func() {
				Expect(err).ToNot(HaveOccurred())
			})
			It("should persist applications successfully", func() {
				// var previusCall *gomock.Call
				// hasPreviusCall := false
				// for _, app := range appConfig.Apps {
				// 	if !hasPreviusCall {
				// 		previusCall = mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
				// 		hasPreviusCall = true
				// 	} else {
				// 		previusCall = mockRepo.EXPECT().Set(app).After(previusCall)
				// 	}
				// }
				appConfig.Repo = mockRepo
				Expect(appConfig.Repo).To(Equal(mockRepo))

				// mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
				err := appConfig.Persists()
				// mockRepo.EXPECT().Set(gomock.Any()).Return(nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
