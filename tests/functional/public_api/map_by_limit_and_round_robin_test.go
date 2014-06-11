package public_api_test

// import (
// 	. "github.com/innotech/hydra/tests/helpers"
// 	. "github.com/innotech/hydra/vendors/github.com/innotech/hydra-worker-lib"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
// 	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"

// 	"bytes"
// 	// "encoding/json"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"
// )

// var _ = Describe("Public API with map by limit and round robin balancers", func() {
// 	pwd, err := os.Getwd()
// 	if err != nil {
// 		Fail("Failure to apply the absolute path to the current directory")
// 		os.Exit(1)
// 	}
// 	var FIXTURES_PATH string = pwd + "/../../fixtures/"

// 	app_config := FIXTURES_PATH + "apps_for_mapbylimit_and_roundrobin.json"
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

// 	// Run MapByLimit Worker
// 	go func() {
// 		mapByLimitWorker := NewWorker("tcp://"+loadBalancerAddr, "MapByLimit", false)
// 		fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
// 			limitAttr := args["limitAttr"].(string)
// 			limitValue, _ := strconv.ParseFloat(args["limitValue"].(string), 64)
// 			mapSort := args["mapSort"].(string)

// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("+ INITIAL INSTANCES +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Printf("%#v", instances)
// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

// 			mappedInstances := make([][]map[string]interface{}, 2)
// 			for _, i := range instances {
// 				instance := i.(map[string]interface{})
// 				// if val, ok := mappedInstances[instance[limitAttr]]; ok {
// 				// if val < limitValue {
// 				value, _ := strconv.ParseFloat(instance["Info"].(map[string]interface{})[limitAttr].(string), 64)
// 				if value < limitValue {
// 					mappedInstances[0] = append(mappedInstances[0], instance)
// 				} else {
// 					mappedInstances[1] = append(mappedInstances[1], instance)
// 				}
// 				// } else {
// 				// 	// TODO: Send an error
// 				// }
// 			}

// 			computedInstances := make([]interface{}, 2)
// 			if mapSort == "reverse" {
// 				computedInstances[0] = mappedInstances[1]
// 				computedInstances[1] = mappedInstances[0]
// 			} else {
// 				// computedInstances = mappedInstances
// 				computedInstances[0] = mappedInstances[0]
// 				computedInstances[1] = mappedInstances[1]
// 			}

// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("+ MAP BY LIMIT	++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Printf("%#v", computedInstances)
// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			return computedInstances
// 		}
// 		mapByLimitWorker.Run(fn)
// 	}()

// 	// Run RoundRobin Worker
// 	// var lastInstanceIndex int = 0
// 	var lastInstanceIndex map[string][]int = make(map[string][]int)

// 	sortSlice := func(instances []interface{}, fisrtElement int, appId string, iteration int) []interface{} {
// 		var index int = lastInstanceIndex[appId][iteration]
// 		computedInstances := make([]interface{}, 0)
// 		// computedInstances = append(computedInstances, instances[lastInstanceIndex+1:])
// 		// computedInstances = append(computedInstances, instances[:lastInstanceIndex+1])
// 		if index < len(instances) {
// 			computedInstances = append(computedInstances, instances[index:])
// 		}
// 		if index > 0 {
// 			computedInstances = append(computedInstances, instances[:index])
// 		}

// 		if index < len(instances)-1 {
// 			lastInstanceIndex[appId][iteration] = lastInstanceIndex[appId][iteration] + 1
// 		} else {
// 			lastInstanceIndex[appId][iteration] = 0
// 		}

// 		return computedInstances
// 	}

// 	go func() {
// 		roundRobinWorker := NewWorker("tcp://"+loadBalancerAddr, "RoundRobin", false)
// 		fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
// 			// var computedInstances []interface{}

// 			// log.Printf("--- LAST INSTANCE INDEX --- %d", lastInstanceIndex)
// 			// if len(instances) > lastInstanceIndex {
// 			// 	if len(instances) > lastInstanceIndex+1 {
// 			// 		computedInstances = sortSlice(instances, lastInstanceIndex+1)
// 			// 	}
// 			// } else {
// 			// 	rand.Seed(time.Now().Unix())
// 			// 	randomIndex := rand.Intn(len(instances))
// 			// 	computedInstances = sortSlice(instances, randomIndex)
// 			// }

// 			// log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			// log.Println("+ ROUND ROBIN +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			// log.Printf("%#v", computedInstances)
// 			// log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			// log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			// return computedInstances

// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n\n")
// 			log.Println("+ ROUND ROBIN start +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("\n\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

// 			log.Printf("Instances: %#v", instances)
// 			log.Printf("Args: %#v", args)

// 			var computedInstances []interface{}

// 			appId := args["appId"].(string)
// 			log.Printf("appId: %s", appId)
// 			iteration := args["iteration"].(int)
// 			log.Printf("iteration: %d", iteration)
// 			if _, ok := lastInstanceIndex[appId]; !ok {
// 				lastInstanceIndex[appId] = make([]int, 0)
// 			}
// 			if iteration >= len(lastInstanceIndex[appId]) {
// 				lastInstanceIndex[appId] = append(lastInstanceIndex[appId], 0)
// 			}
// 			var index int = lastInstanceIndex[appId][iteration]
// 			log.Printf("--- LAST INSTANCE INDEX --- %d", index)
// 			if len(instances) > index {
// 				log.Println("Entra 1")
// 				// if len(instances) > index+1 {
// 				// 	log.Println("Entra 3")
// 				computedInstances = sortSlice(instances, index+1, appId, iteration)
// 				// }
// 			} else {
// 				log.Println("Entra 2")
// 				rand.Seed(time.Now().Unix())
// 				randomIndex := rand.Intn(len(instances))
// 				computedInstances = sortSlice(instances, randomIndex, appId, iteration)
// 			}

// 			log.Println("+ ROUND ROBIN computedInstances +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Printf("%#v", computedInstances)

// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++\n\n")
// 			log.Println("+ ROUND ROBIN end +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 			log.Println("\n\n+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

// 			return computedInstances
// 		}
// 		roundRobinWorker.Run(fn)
// 	}()

// 	// Fake Hydra Probe
// 	RunFakeHydraProbe := func(instanceInfo []byte) {
// 		for {
// 			httpUtils.Post(app1InstancesAddr, "application/json", bytes.NewReader(instanceInfo))
// 			// log.Println("++++++++++++++++++++++++ Fake Hydra Probe EMIT +++++++++++++++++++++++")
// 			time.Sleep(5 * time.Second)
// 		}
// 	}

// 	// Run Probes
// 	go RunFakeHydraProbe([]byte(`{"PC1101": {"uri": "http://pc1101:8080", "limit": 25.4}}`))
// 	time.Sleep(100 * time.Millisecond)
// 	go RunFakeHydraProbe([]byte(`{"PC1102": {"uri": "http://pc1102:8080", "limit": 75.7}}`))
// 	time.Sleep(100 * time.Millisecond)
// 	go RunFakeHydraProbe([]byte(`{"PC1103": {"uri": "http://pc1103:8080", "limit": 13.0}}`))
// 	time.Sleep(100 * time.Millisecond)
// 	go RunFakeHydraProbe([]byte(`{"PC1104": {"uri": "http://pc1104:8080", "limit": 91.3}}`))
// 	time.Sleep(100 * time.Millisecond)
// 	go RunFakeHydraProbe([]byte(`{"PC1105": {"uri": "http://pc1105:8080", "limit": 45.9}}`))

// 	time.Sleep(2 * time.Second)
// 	var response *http.Response
// 	var errG error
// 	go func() {
// 		log.Println("---------- SEND GET 1 ----------")
// 		response, err = httpUtils.Get("http://" + publicAddr + "/apps/App1")
// 		log.Println("---------- END SEND GET 1 ----------")
// 		log.Printf("%#v", response)
// 		log.Printf("%#v", response.Body)
// 		uriList, _ := httpUtils.ReadBodyStringArray(response)
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Println("+ URI 1 OOOOOOOOOO ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Printf("%#v", uriList)
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	}()
// 	time.Sleep(12 * time.Second)

// 	// log.Println("---------- START TESTS ----------")
// 	// Context("When Ping application exist and Pong worker is registered", func() {
// 	// 	Describe("Sending a correct request to get Ping uris", func() {
// 	// 		It("should receive the correct list of URIs", func() {
// 	// 			// log.Println("---------- START EXPECTS ----------")
// 	// 			Expect(errG).NotTo(HaveOccurred())
// 	// 			Expect(response.StatusCode).To(Equal(200))
// 	// 			uriList, err2 := httpUtils.ReadBodyStringArray(response)
// 	// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	// 			log.Println("+ URI LIST ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	// 			log.Printf("%#v", uriList)
// 	// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	// 			log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	// 			Expect(err2).NotTo(HaveOccurred())
// 	// 			Expect(uriList).To(HaveLen(5))
// 	// 			// log.Println("Processed response %#v", uriList)
// 	// 			// Expect(uriList[0]).To(Equal("http://pc1104:8080"))
// 	// 			// Expect(uriList[1]).To(Equal("http://pc1102:8080"))
// 	// 			// Expect(uriList[2]).To(Equal("http://pc1103:8080"))
// 	// 			// Expect(uriList[3]).To(Equal("http://pc1105:8080"))
// 	// 			// Expect(uriList[4]).To(Equal("http://pc1101:8080"))

// 	// 			Expect(uriList[0]).To(Equal("http://pc1102:8080"))
// 	// 			Expect(uriList[1]).To(Equal("http://pc1104:8080"))
// 	// 			Expect(uriList[2]).To(Equal("http://pc1101:8080"))
// 	// 			Expect(uriList[3]).To(Equal("http://pc1103:8080"))
// 	// 			Expect(uriList[4]).To(Equal("http://pc1105:8080"))
// 	// 		})
// 	// 	})
// 	// })

// 	// time.Sleep(2 * time.Second)
// 	// var response *http.Response
// 	// var errG error
// 	go func() {
// 		log.Println("---------- SEND GET 2 ----------")
// 		response, err = httpUtils.Get("http://" + publicAddr + "/apps/App1")
// 		log.Println("---------- END SEND GET 2 ----------")
// 		log.Printf("%#v", response)
// 		log.Printf("%#v", response.Body)
// 		uriList, _ := httpUtils.ReadBodyStringArray(response)
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Println("+ URI 2 OOOOOOOOOO ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Printf("%#v", uriList)
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 		log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 	}()
// 	time.Sleep(2 * time.Second)

// 	log.Println("---------- START TESTS 2 ----------")
// 	Context("When Ping application exist and Pong worker is registered", func() {
// 		Describe("Sending a correct request to get Ping uris", func() {
// 			It("should receive the correct list of URIs", func() {
// 				// log.Println("---------- START EXPECTS ----------")
// 				Expect(errG).NotTo(HaveOccurred())
// 				Expect(response.StatusCode).To(Equal(200))
// 				uriList, err2 := httpUtils.ReadBodyStringArray(response)
// 				log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 				log.Println("+ URI LIST 2 ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 				log.Printf("%#v", uriList)
// 				log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 				log.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
// 				Expect(err2).NotTo(HaveOccurred())
// 				Expect(uriList).To(HaveLen(5))
// 				// log.Println("Processed response %#v", uriList)
// 				// Expect(uriList[0]).To(Equal("http://pc1102:8080"))
// 				// Expect(uriList[1]).To(Equal("http://pc1104:8080"))
// 				// Expect(uriList[2]).To(Equal("http://pc1105:8080"))
// 				// Expect(uriList[3]).To(Equal("http://pc1101:8080"))
// 				// Expect(uriList[4]).To(Equal("http://pc1103:8080"))

// 				Expect(uriList[0]).To(Equal("http://pc1104:8080"))
// 				Expect(uriList[1]).To(Equal("http://pc1102:8080"))
// 				Expect(uriList[2]).To(Equal("http://pc1103:8080"))
// 				Expect(uriList[3]).To(Equal("http://pc1105:8080"))
// 				Expect(uriList[4]).To(Equal("http://pc1101:8080"))
// 			})
// 		})
// 	})
// })
