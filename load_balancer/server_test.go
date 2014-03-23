package load_balancer_test

import (
	. "github.com/innotech/hydra/load_balancer"
	. "github.com/innotech/hydra/vendors/github.com/onsi/ginkgo"
	. "github.com/innotech/hydra/vendors/github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Describe("Processing Worker message", func() {
		Context("When message doesn't contain command", func() {
			bs := NewLoadBalancer("ipc://frontend.ipc", "tcp://localhost:7777")
			defer bs.Close()
			sender := []byte("worker-1")
			msg := [][]byte{}
			err := bs.ProcessWorker(sender, msg)
			It("should throw an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("When the command is READY", func() {
			bs := NewLoadBalancer("ipc://frontend.ipc", "tcp://localhost:7777")
			defer bs.Close()
			sender := []byte("worker-1")
			identity := hex.EncodeToString(sender)
			Context("When message doesn't contain service name", func() {
				msg := [][]byte{[]byte(SIGNAL_READY)}
				err := bs.ProcessWorker(sender, msg)
				It("should throw an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
	// Describe("Extracting app ID from multipart message", func() {
	// 	s := NewBalancerServer()
	// 	Context("When multipart message doesn't contain application ID", func() {
	// 		parts := make([][]byte, 1)
	// 		parts[0] = []byte("client-identity")
	// 		appId, err := s.ExtractAppIDFromClientMultipartMsg(parts)
	// 		It("should throw an error", func() {
	// 			Expect(err).To(HaveOccurred())
	// 		})
	// 	})
	// 	Context("When multipart message contains an empty application ID", func() {
	// 		parts := make([][]byte, 2)
	// 		parts[0] = []byte("client-identity")
	// 		parts[1] = []byte{}
	// 		appId, err := s.ExtractAppIDFromClientMultipartMsg(parts)
	// 		It("should throw an error", func() {
	// 			Expect(err).To(HaveOccurred())
	// 		})
	// 	})
	// 	Context("When multipart message is not correct", func() {
	// 		const APP_ID string = "app-1"
	// 		parts := make([][]byte, 2)
	// 		parts[0] = []byte("client-identity")
	// 		parts[1] = []byte(APP_ID)
	// 		appId, err := s.ExtractAppIDFromClientMultipartMsg(parts)
	// 		It("should throw an error", func() {
	// 			Expect(err).NotTo(HaveOccurred())
	// 			Expect(appId).To(Equal(APP_ID))
	// 		})
	// 	})
	// })
	// Describe("Extracting balancer pipeline from application data", func() {
	// 	Context("When data doesn't contain balancers key", func() {
	// 		data := map[string]interface{} {
	// 			"intances": make([]string),
	// 		}
	// 		b := NewBalancerServer()
	// 		pipeline, err := b.ExtractBalancerPipelineFromApplicationData()
	// 		It("should throw an error", func() {
	// 			Expect(err).To(HaveOccurred())
	// 		})
	// 	})
	// 	Context("When balancers is empty", func() {
	// 		data := map[string]interface{} {
	// 			"balancers": make(map[string]map[string]interface{}),
	// 			"intances": make([]string),
	// 		}
	// 		b := NewBalancerServer()
	// 		pipeline, err := b.ExtractBalancerPipelineFromApplicationData()
	// 		It("should throw an error", func() {
	// 			Expect(err).To(HaveOccurred())
	// 		})
	// 	})
	// 	Context("When balancers contains balancers", func() {
	// 		data := map[string]interface{} {
	// 			"balancers": map[string]map[string]interface{} {
	// 				"balancer-1": map[string]interface{} {
	// 					"b1-attr1": "value1",
	// 					"b1-attr2": 2,
	// 				},
	// 				"balancer-2": map[string]interface{} {
	// 					"b2-attr1": "value2",
	// 					"b2-attr2": 4,
	// 				},
	// 			},
	// 			"intances": make([]string),
	// 		}
	// 		b := NewBalancerServer()
	// 		pipeline, err := b.ExtractBalancerPipelineFromApplicationData()
	// 		It("should be return an correct pipeline", func() {
	// 			Expect(err).ToNot(HaveOccurred())
	// 			Expect(pipeline).ToNot(BeNil())
	// 			Expect(pipeline).To(HaveLen(2))
	// 			// Expect(pipeline[])
	// 		})
	// 	})
	// })

	// Describe("Making new instance", func() {
	// 	conf := config.Balancer{
	// 		addr: "*:5555",
	// 	}
	// 	s := NewBalancerServer(conf)
	// 	It("should set the attributes correctly", func() {
	// 		Expect(s.Addr).To(Equal("*:5555"))
	// 	})
	// })
	// Describe("Receiving a new request", func() {
	// 	requestedApplication := "app-1"
	// 	s := NewBalancerServer(conf)
	// 	Describe("Getting aplication from database", func() {
	// 		app := s.AppRepository.Get("bad-"+requestedApplication)
	// 		Context("When application doesn't exist", func() {
	// 			// TODO
	// 		})
	// 		app := s.AppRepository.Get(requestedApplication)
	// 		Context("When application exists", func() {
	// 			// TODO
	// 			Fail("Failure reason")
	// 		})
	// 	})
	// })
	// Describe("Making new pipeline", func() {
	// 	// TODO new etcd base model
	// 	p := NewPipeline(balancers)

	// })
	// Describe("Registering application plumber", func() {
	// 	b := NewServer()
	// 	appId, appAttrs :=
	// 	b.RegisterPlumber(app)
	// 	It("should be registered succesfully", func() {
	// 		Expect(b, ...)

	// 	})
	// })
})
