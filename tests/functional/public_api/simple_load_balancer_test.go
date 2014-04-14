package public_api_test

import (
	. "github.com/innotech/hydra/tests/helpers"
	. "github.com/innotech/hydra/vendors/github.com/innotech/hydra-worker-lib"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	// "encoding/json"
	"log"
	"net/http"
	"os"
	// "strings"
	"time"
)

var _ = Describe("Public API", func() {
	pwd, err := os.Getwd()
	if err != nil {
		Fail("Failure to apply the absolute path to the current directory")
		os.Exit(1)
	}
	var FIXTURES_PATH string = pwd + "/../../fixtures/"

	app_config := FIXTURES_PATH + "apps_ping.json"
	hydra_name := "node0"
	data_dir_path := DATA_DIR_PATH + hydra_name
	loadBalancerAddr := "127.0.0.1:7777"
	privateAddr := "127.0.0.1:7771"
	publicAddr := "127.0.0.1:7772"
	pingInstancesAddr := "http://" + privateAddr + "/apps/Ping/Instances"

	httpUtils := NewHTTPClientHelper()

	args := []string{"-name=" + hydra_name, "-load-balancer-addr=" + loadBalancerAddr, "-private-addr=" + privateAddr, "-public-addr=" + publicAddr, "-data-dir=" + data_dir_path, "-apps-file=" + app_config}
	process := RunHydraInStandaloneAndReturnProcess(args)
	defer KillHydraProcess(process)
	time.Sleep(5 * time.Second)

	// Run Pong Worker
	go func() {
		pongWorker := NewWorker("tcp://"+loadBalancerAddr, "Pong", false)
		fn := func(instances []map[string]interface{}, args map[string]string) []interface{} {
			computedInstances := make([]interface{}, 0)
			for _, instance := range instances {
				computedInstances = append(computedInstances, instance["Info"].(map[string]interface{})["uri"].(string))
			}
			return computedInstances
		}
		log.Println("++++++++++++++++++++++++ PRE pongWorker RUN +++++++++++++++++++++++")
		pongWorker.Run(fn)
	}()

	// Run Fake Hydra Probe
	go func() {
		instanceInfo := []byte(`{"PC1148": {"memLoad": 27.5, "uri": "ssh://localhost:22", "connections": 2, "cpuLoad": 15.4, "state": 0, "cost": "5", "cloud": "susecloud"}}`)
		// instanceMap := map[string]interface{}{
		// 	"PC1148": map[string]interface{}{
		// 		"uri": "ssh://localhost:22",
		// 	},
		// }
		// instanceInfo, _ := json.Marshal(instanceMap)
		for {
			httpUtils.Post(pingInstancesAddr, "application/json", bytes.NewReader(instanceInfo))
			log.Println("++++++++++++++++++++++++ Fake Hydra Probe EMIT +++++++++++++++++++++++")
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)
	var response *http.Response
	var errG error
	go func() {
		log.Println("---------- SEND GET 1 ----------")
		response, err = httpUtils.Get("http://" + publicAddr + "/apps/Ping")
		log.Println("---------- END SEND GET 1 ----------")
		log.Printf("%#v", response)
		log.Printf("%#v", response.Body)
	}()
	time.Sleep(12 * time.Second)

	log.Println("---------- START TESTS ----------")
	// test := func(chan bool) {
	Context("When Ping application exist and Pong worker is registered", func() {
		Describe("Sending a correct request to get Ping uris", func() {
			// log.Println("---------- SEND GET ----------")
			// response, err := httpUtils.Get("http://" + publicAddr + "/apps/Ping")
			// log.Printf("%#v", response)
			It("should receive the correct list of URIs", func() {
				// var response *http.Response
				// var err error
				// c := make(chan bool)
				// go func() {
				// 	log.Println("---------- SEND GET 2 ----------")
				// 	response, err = httpUtils.Get("http://" + publicAddr + "/apps/Ping")
				// 	log.Printf("%#v", response)
				// 	c <- true
				// }()
				// <-c
				log.Println("---------- START EXPECTS ----------")
				Expect(errG).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				// var uriList []string
				uriList, err2 := httpUtils.ReadBodyStringArray(response)
				// err := json.Unmarshal(response, uriList)
				Expect(err2).NotTo(HaveOccurred())
				Expect(uriList).To(HaveLen(1))
				log.Println("Processed response %#v", uriList)
				// Expect(uriList[0]).To(Equal("ssh://localhost:22"))
			})
		})
	})
	// }

	// time.Sleep(12 * time.Second)
})
