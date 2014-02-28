package tests_test

import (
	. "github.com/innotech/hydra/tests/functional"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	// "testing"
)

var _ = Describe("SigleNode", func() {
	// Describe("Starting as standalone", func() {
	procAttr := new(os.ProcAttr)
	// procAttr := new(os.ProcAttr)
	// procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	args := []string{"hydra", "-name=node1", "-f", "-data-dir=/tmp/node1"}

	process, err := os.StartProcess(HydraBinPath, args, procAttr)
	// process, err := os.StartProcess(HydraBinPath, args, nil)
	It("should be running successfully", func() {
		Expect(err).NotTo(HaveOccurred())
	})
	if err != nil {
		GinkgoT().Fatal("start process failed:" + err.Error())
		return
	}
	defer process.Kill()

	time.Sleep(2 * time.Second)
	Describe("Setting an application", func() {
		Context("When PUT request is defined correctly", func() {
			client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
			// appJson := `{
			// 	“localStrategyEvents”: {...},
			// 	“cloudStrategyEvents”: {...},
			// 	"servers":[
			// 	    {
			// 			"server":"http://mycompany.com/api",
			// 			"status":"statusStrut",
			// 			“cost”: 3,
			// 			“cloud”: “amazon”
			// 		}
			// 	]
			// }`

			// appJson := `{
			// 	"localStrategyEvents": {...},
			// 	"cloudStrategyEvents": {...},
			// 	"servers": [
			// 	    {
			// 			"server": "http://mycompany.com/api",
			// 			"status": {
			//  						cpuLoad: 30,
			//  						memLoad: 50,
			//  						timeStamp: 42374897239
			//  					},
			// 			"cost": 3,
			// 			"cloud": "amazon"
			// 		}
			// 	]
			// }`

			appJson := `{
				"server": "http://mycompany.com/api",
				"status": {
					cpuLoad: 30,
					memLoad: 50,
					timeStamp: 42374897239
				},
				"cost": 3,
				"cloud": "amazon"
			}`

			b := strings.NewReader(appJson)
			req, err := http.NewRequest("PUT", "http://127.0.0.1:8082/applications/testapp/instances/mycompany", b)
			if err != nil {
				fmt.Println("Bad request 1")
				return
			}
			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if res == nil {
				fmt.Println("Bad request 2")
				// return
			}
			It("should be received a correct response", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
	// })
})
