package public_api_test

// import (
// 	. "github.com/innotech/hydra/tests/helpers"
// 	. "github.com/innotech/hydra/vendors/github.com/innotech/hydra-worker-lib"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	"bytes"
// 	"net/http"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"time"
// )

// const (
// 	decr string = "0"
// 	incr string = "1"
// )

// var order, sortAttr string

// type Instances []map[string]interface{}

// func (a Instances) Len() int      { return len(a) }
// func (a Instances) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
// func (s Instances) Less(i, j int) bool {
// 	var less bool
// 	if order == decr {
// 		a1, _ := strconv.ParseFloat(s[i]["Info"].(map[string]interface{})[sortAttr].(string), 64)
// 		a2, _ := strconv.ParseFloat(s[j]["Info"].(map[string]interface{})[sortAttr].(string), 64)
// 		less = a1 > a2
// 	} else {
// 		a1, _ := strconv.ParseFloat(s[i]["Info"].(map[string]interface{})[sortAttr].(string), 64)
// 		a2, _ := strconv.ParseFloat(s[j]["Info"].(map[string]interface{})[sortAttr].(string), 64)
// 		less = a1 < a2
// 	}
// 	return less
// }

// var _ = Describe("Public API with double balancer", func() {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		Fail("Failure to apply the absolute path to the current directory")
// 		os.Exit(1)
// 	}
// 	var FIXTURES_PATH string = pwd + "/../../fixtures/"

// 	app_config := FIXTURES_PATH + "apps_for_double_balancer.json"
// 	hydra_name := "node0"
// 	data_dir_path := DATA_DIR_PATH + hydra_name
// 	loadBalancerAddr := "127.0.0.1:7777"
// 	privateAddr := "127.0.0.1:7771"
// 	publicAddr := "127.0.0.1:7772"
// 	app1InstancesAddr := "http://" + privateAddr + "/apps/App1/Instances"

// 	httpUtils := NewHTTPClientHelper()

// 	args := []string{"-name=" + hydra_name, "-load-balancer-addr=" + loadBalancerAddr, "-private-addr=" + privateAddr, "-public-addr=" + publicAddr, "-data-dir=" + data_dir_path, "-apps-file=" + app_config}
// 	process := RunHydraInStandaloneAndReturnProcess(args)
// 	defer KillHydraProcess(process)
// 	time.Sleep(5 * time.Second)

// 	// Run MapAndSort Worker
// 	go func() {
// 		mapAndSortWorker := NewWorker("tcp://"+loadBalancerAddr, "MapAndSort", false)
// 		fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
// 			var mappedInstances map[string][]interface{}
// 			mappedInstances = make(map[string][]interface{})
// 			for _, i := range instances {
// 				instance := i.(map[string]interface{})
// 				r := instance["Info"].(map[string]interface{})[args["mapAttr"].(string)].(string)
// 				if len(mappedInstances[r]) == 0 {
// 					mappedInstances[r] = make([]interface{}, 0)
// 				}
// 				if _, ok := mappedInstances[r]; ok {
// 					mappedInstances[r] = append(mappedInstances[r], instance)
// 				} else {
// 					mappedInstances[r] = []interface{}{instance}
// 				}
// 			}

// 			computedInstances := make([]interface{}, 0)
// 			for _, mapAttr := range args["mapSort"].(map[string]interface{}) {
// 				computedInstances = append(computedInstances, mappedInstances[mapAttr.(string)])
// 			}

// 			return computedInstances
// 		}
// 		mapAndSortWorker.Run(fn)
// 	}()

// 	// Run SortByNumber Worker
// 	go func() {
// 		const (
// 			DECR string = "0"
// 			INCR string = "1"
// 		)

// 		sortByNumberWorker := NewWorker("tcp://"+loadBalancerAddr, "SortByNumber", false)
// 		fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
// 			var finalInstances []map[string]interface{}
// 			finalInstances = make([]map[string]interface{}, 0)
// 			for _, instance := range instances {
// 				finalInstances = append(finalInstances, instance.(map[string]interface{}))
// 			}

// 			sortAttr = args["sortAttr"].(string)
// 			order = args["order"].(string)
// 			sort.Sort(Instances(finalInstances))

// 			var finalInstances2 []interface{}
// 			finalInstances2 = make([]interface{}, 0)
// 			for _, instance := range finalInstances {
// 				finalInstances2 = append(finalInstances2, instance)
// 			}

// 			return finalInstances2
// 		}
// 		sortByNumberWorker.Run(fn)
// 	}()

// 	// Fake Hydra Probe
// 	RunFakeHydraProbe := func(instanceInfo []byte) {
// 		for {
// 			httpUtils.Post(app1InstancesAddr, "application/json", bytes.NewReader(instanceInfo))
// 			time.Sleep(5 * time.Second)
// 		}
// 	}

// 	// Run Probes
// 	go RunFakeHydraProbe([]byte(`{"PC1001": {"uri": "http://pc1001:8080", "cpuLoad": 15.4, "cloud": "amazon"}}`))
// 	go RunFakeHydraProbe([]byte(`{"PC1002": {"uri": "http://pc1002:8080", "cpuLoad": 45.7, "cloud": "azure"}}`))
// 	go RunFakeHydraProbe([]byte(`{"PC1003": {"uri": "http://pc1003:8080", "cpuLoad": 19.0, "cloud": "amazon"}}`))
// 	go RunFakeHydraProbe([]byte(`{"PC1004": {"uri": "http://pc1004:8080", "cpuLoad": 11.3, "cloud": "azure"}}`))
// 	go RunFakeHydraProbe([]byte(`{"PC1005": {"uri": "http://pc1005:8080", "cpuLoad": 85.9, "cloud": "google"}}`))

// 	time.Sleep(2 * time.Second)
// 	var response *http.Response
// 	var errG error
// 	go func() {
// 		response, err = httpUtils.Get("http://" + publicAddr + "/apps/App1")
// 	}()
// 	time.Sleep(12 * time.Second)

// 	Context("When Ping application exist and Pong worker is registered", func() {
// 		Describe("Sending a correct request to get Ping uris", func() {
// 			It("should receive the correct list of URIs", func() {
// 				Expect(errG).NotTo(HaveOccurred())
// 				Expect(response.StatusCode).To(Equal(200))
// 				uriList, err2 := httpUtils.ReadBodyStringArray(response)
// 				Expect(err2).NotTo(HaveOccurred())
// 				Expect(uriList).To(HaveLen(5))
// 				Expect(uriList[0]).To(Equal("http://pc1005:8080"))
// 				Expect(uriList[1]).To(Equal("http://pc1001:8080"))
// 				Expect(uriList[2]).To(Equal("http://pc1003:8080"))
// 				Expect(uriList[3]).To(Equal("http://pc1004:8080"))
// 				Expect(uriList[4]).To(Equal("http://pc1002:8080"))
// 			})
// 		})
// 	})
// })
