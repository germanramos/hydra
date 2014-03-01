package utils_test

import (
	. "github.com/innotech/hydra/utils"
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
				ThridLevelString string
				// ThirdLevelArray  []string
				ThirdLevelObject struct {
					FourthLevelInt int
				}
			}
			SecondLevelInt    int
			SecondLevelFloat  float64
			SecondLevelBool   bool
			SecondLevelString string
			// SecondLevelArray []string
		}
		FirstLevelInt          int
		FirstLevelFloat        float64
		FirstLevelString       string
		FirstLevelBool         bool
		FirstLevelStringsArray []string
		FirstLevelMap          map[string]interface{}
	}

	var jsonBlob = []byte(`{
		"FirstLevelObject": {
			"SecondLevelInt": 2345,
			"SecondLevelFloat": 5.4,
			"SecondLevelBool": false,
			"SecondLevelString": "second string test",
			"SecondLevelObject": {
			}
		},
		"FirstLevelInt": 1234,
		"FirstLevelFloat": 10.234,
		"FirstLevelBool": true,
		"FirstLevelString": "string test",
		"FirstLevelStringsArray": ["string1", "string2", "string 1 + 2"],
		"FirstLevelMap": {
			"MapSecondLevelInt": 5,
			"MapSecondLevelFloat": 40.44,
			"MapSecondLevelBool": false,
			"MapSecondLevelMapsArray": [{
				"test1": 1	
			}, {
				"test2": "fred"
			}],
			"MapSecondLevelMap": {
				"MapThirdLevelString": "Yuhuu!"
			}
		}
	}`)

	// var jsonBlob = []byte(`{
	// 	"FirstLevelObject": {
	// 		"SecondLevelInt": 2345,
	// 		"SecondLevelFloat": 5.4,
	// 		"SecondLevelBool": false,
	// 		"SecondLevelString": "second string test",
	// 		"SecondLevelObject": {
	// 		}
	// 	},
	// 	"FirstLevelInt": 1234,
	// 	"FirstLevelFloat": 10.234,
	// 	"FirstLevelBool": true,
	// 	"FirstLevelString": "string test"
	// }`)

	Describe("decoding json object", func() {
		// When whole equal true
		FContext("when the json object is correct", func() {
			var cyclicStruct json
			ops, err := DecodeJsonObject(jsonBlob, &cyclicStruct)
			It("should be decoded successfully", func() {
				Expect(err).To(BeNil(), "error should be nil")
				Expect(ops["/FirstLevelInt"]).To(Equal("1234"), `ops["/FirstLevelInt"] should be equal 1234`)
				Expect(ops["/FirstLevelFloat"]).To(Equal("10.23"), `ops["/FirstLevelFloat"] should be equal 10.23`)
				Expect(ops["/FirstLevelBool"]).To(Equal("true"), `ops["/FirstLevelFloat"] should be equal true`)
				Expect(ops["/FirstLevelString"]).To(Equal("string test"), `ops["/FirstLevelString"] should be equal "string test"`)
				Expect(ops["/FirstLevelStringsArray/0"]).To(Equal("string1"), `ops["/FirstLevelStringsArray/0"] should be equal "string1"`)
				Expect(ops["/FirstLevelStringsArray/1"]).To(Equal("string2"), `ops["/FirstLevelStringsArray/0"] should be equal "string2"`)
				Expect(ops["/FirstLevelStringsArray/2"]).To(Equal("string 1 + 2"), `ops["/FirstLevelStringsArray/0"] should be equal "string 1 + 2"`)
				Expect(ops["/FirstLevelMap/MapSecondLevelInt"]).To(Equal("5.00"), `ops["/FirstLevelMap/MapSecondLevelInt"] should be equal 5.00`)
				Expect(ops["/FirstLevelMap/MapSecondLevelFloat"]).To(Equal("40.44"), `ops["/FirstLevelMap/MapSecondLevelFloat"] should be equal 40.44`)
				Expect(ops["/FirstLevelMap/MapSecondLevelBool"]).To(Equal("false"), `ops["/FirstLevelMap/MapSecondLevelBool"] should be equal false`)
				Expect(ops["/FirstLevelMap/MapSecondLevelMapsArray/0/test1"]).To(Equal("1.00"), `ops["/FirstLevelMap/MapSecondLevelMapsArray/0/test1"] should be equal 1.00`)
				Expect(ops["/FirstLevelMap/MapSecondLevelMapsArray/1/test2"]).To(Equal("fred"), `ops["/FirstLevelMap/MapSecondLevelMapsArray/1/test2"] should be equal "fred"`)
				Expect(ops["/FirstLevelMap/MapSecondLevelMap/MapThirdLevelString"]).To(Equal("Yuhuu!"), `ops["/FirstLevelMap/MapSecondLevelMap/MapThirdLevelString"] should be equal "Yuhuu!"`)
				Expect(ops["/FirstLevelObject/SecondLevelInt"]).To(Equal("2345"), `ops["/FirstLevelObject/SecondLevelInt"] should be equal 2345`)
				Expect(ops["/FirstLevelObject/SecondLevelFloat"]).To(Equal("5.40"), `ops["/FirstLevelObject/SecondLevelFloat"] should be equal 5.43`)
				Expect(ops["/FirstLevelObject/SecondLevelBool"]).To(Equal("false"), `ops["/FirstLevelObject/SecondLevelBool"] should be equal false`)
				Expect(ops["/FirstLevelObject/SecondLevelString"]).To(Equal("second string test"), `ops["/FirstLevelObject/SecondLevelString"] should be equal "second string test"`)
				Expect(ops["/FirstLevelObject/SecondLevelObject/ThirdLevelInt"]).To(Equal("0"))
				Expect(ops["/FirstLevelObject/SecondLevelObject/ThirdLevelFloat"]).To(Equal("0.00"))
				Expect(ops["/FirstLevelObject/SecondLevelObject/ThirdLevelBool"]).To(Equal("false"))
				Expect(ops["/FirstLevelObject/SecondLevelObject/ThirdLevelString"]).To(Equal(""))
				Expect(ops["/FirstLevelObject/SecondLevelObject/ThirdLevelObject/FourthLevelInt"]).To(Equal("0"))
			})
		})
		// Context("when the json object is incorrect", func() {
		// 	var cyclicStruct json
		// 	ops, err := DecodeJsonObject(jsonBlob, &cyclicStruct)
		// 	It("should be decoded successfully", func() {
		// 		Expect(err).To(BeNil(), "error should be nil")
		// 		Expect(ops["FirstLevelInt"]).To(Equal("1234"), `ops["FirstLevelInt"] should be equal 1234`)
		// 	})
		// })
	})
})
