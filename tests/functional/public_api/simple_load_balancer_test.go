package public_api_test

import (
	. "github.com/innotech/hydra/tests/helpers"
	. "github.com/innotech/hydra/vendors/github.com/innotech/hydra-worker-lib"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

	"bytes"
	"log"
	"net/http"
	"os"
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
		fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
			return instances
		}
		pongWorker.Run(fn)
	}()

	// Run Fake Hydra Probe
	go func() {
		instanceInfo := []byte(`{"PC1148": {"memLoad": 27.5, "uri": "ssh://localhost:22", "connections": 2, "cpuLoad": 15.4, "state": 0, "cost": "5", "cloud": "susecloud"}}`)
		for {
			httpUtils.Post(pingInstancesAddr, "application/json", bytes.NewReader(instanceInfo))
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)
	var response *http.Response
	var errG error
	go func() {
		response, err = httpUtils.Get("http://" + publicAddr + "/apps/Ping")
	}()
	time.Sleep(12 * time.Second)

	Context("When Ping application exist and Pong worker is registered", func() {
		Describe("Sending a correct request to get Ping uris", func() {
			It("should receive the correct list of URIs", func() {
				Expect(errG).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))
				uriList, err2 := httpUtils.ReadBodyStringArray(response)
				Expect(err2).NotTo(HaveOccurred())
				Expect(uriList).To(HaveLen(1))
				Expect(uriList[0]).To(Equal("ssh://localhost:22"))
			})
		})
	})
})
