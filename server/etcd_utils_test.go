package server_test

import (
	. "github.com/innotech/hydra/server"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("EtcdUtils", func() {
	type json struct {
		FirstLevelObject struct {
			SecondLevelObject struct {
				ThirdLevelInt    int
				ThirdLevelFloat  float64
				ThirdLevelBool   bool
				ThirdLevelArray  []string
				ThirdLevelObject struct {
					FourthLevelInt int
				}
			}
			SecondLevelInt   int
			SecondLevelFloat float64
			SecondLevelBool  bool
			SecondLevelArray []string
		}
		FirstLevelInt   int
		FirstLevelFloat float64
		FirstLevelBool  bool
		FirstLevelArray []string
	}

	// var jsonBlob = []byte(`{
	// 	"FirstLevelObject": {
	// 		"SecondLevelObject": {
	// 			"ThirdLevelInt": 5467
	// 		}
	// 		"SecondLevelInt": 2345
	// 		"SecondLevelFloat": 5.432
	// 		"SecondLevelBool": false
	// 		"SecondLevelArray": ["home", "car"]
	// 	}
	// 	"FirstLevelInt": 1234
	// 	"FirstLevelFloat": 10.234
	// 	"FirstLevelBool": true
	// 	"FirstLevelArray": ["home", "car"]
	// }`)
	var jsonBlob = []byte(`{
		"FirstLevelInt": 1234
	}`)

	Describe("decoding json object", func() {
		FContext("when the json object is correct", func() {
			var cyclicStruct json
			ops, err := DecodeJsonObject(jsonBlob, &cyclicStruct)
			It("should be decoded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
				Expect(ops["FirstLevelInt"]).To(Equal("1234"), `ops["FirstLevelInt"] should be equal 1234`)
			})
		})
		Context("when the json object is incorrect", func() {
			var cyclicStruct json
			ops, err := DecodeJsonObject(jsonBlob, &cyclicStruct)
			It("should be decoded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
				Expect(ops["FirstLevelInt"]).To(Equal("1234"), `ops["FirstLevelInt"] should be equal 1234`)
			})
		})
	})
})
