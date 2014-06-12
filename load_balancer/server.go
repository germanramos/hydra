package load_balancer

import (
	"encoding/hex"
	"encoding/json"
	"reflect"
	"time"

	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"

	"github.com/innotech/hydra/log"
	. "github.com/innotech/hydra/model/entity"
)

const (
	INTERNAL_SERVICE_PREFIX = "isb."
	// Merge all heartbeat
	HEARTBEAT_INTERVAL = 2500 * time.Millisecond
	HEARTBEAT_EXPIRY   = HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS
)

type Broker interface {
	Close()
	Run()
}

type lbWorker struct {
	identity string     //  UUID Identity of worker
	address  []byte     //  Address to route to
	expiry   time.Time  //  Expires at this point, unless heartbeat
	service  *lbService //  Owning service, if known
}

type lbChain struct {
	app      string
	client   string
	msg      []byte
	shackles *ZList
}

type lbShackle struct {
	serviceName string
	serviceArgs map[string]interface{}
}

type lbService struct {
	broker   Broker
	name     string
	requests [][][]byte //  List of client requests
	waiting  *ZList     //  List of waiting workers
}

type loadBalancer struct {
	context     *zmq.Context          //  Context
	heartbeatAt time.Time             //  When to send HEARTBEAT
	services    map[string]*lbService //  Known services
	frontend    *zmq.Socket           //  Socket for clients
	backend     *zmq.Socket           //  Socket for workers
	waiting     *ZList                //  Idle workers
	workers     map[string]*lbWorker  //  Known workers
	chains      map[string]lbChain    //  Peding Requests
}

func NewLoadBalancer(frontendEndpoint, backendEndpoint string) *loadBalancer {
	context, _ := zmq.NewContext()
	// Define inproc socket to talk with HTTP CLient API
	frontend, _ := context.NewSocket(zmq.ROUTER)
	frontend.SetLinger(0)
	frontend.Bind(frontendEndpoint)
	// log.Infof("Load Balancer Frontend is active at %s\n", frontendEndpoint)
	// Define tcp socket to talk with Workers
	backend, _ := context.NewSocket(zmq.ROUTER)
	backend.SetLinger(0)
	backend.Bind(backendEndpoint)
	// log.Infof("Load Balancer Backend is active at %s\n", backendEndpoint)
	return &loadBalancer{
		context:     context,
		heartbeatAt: time.Now().Add(HEARTBEAT_INTERVAL),
		services:    make(map[string]*lbService),
		frontend:    frontend,
		backend:     backend,
		waiting:     NewList(),
		workers:     make(map[string]*lbWorker),
		chains:      make(map[string]lbChain),
	}
}

// Deletes worker from all data structures, and deletes worker.
func (self *loadBalancer) deleteWorker(worker *lbWorker, disconnect bool) {
	log.Info("----------------------------- Entra en deleteWorker")
	if worker == nil {
		log.Warn("Nil worker")
	}

	if disconnect {
		log.Info("SIGNAL_DISCONNECT")
		self.sendToWorker(worker, SIGNAL_DISCONNECT, nil, nil)
	}

	if worker.service != nil {
		log.Info("Delete Worker Service")
		worker.service.waiting.Delete(worker)
	}
	self.waiting.Delete(worker)
	delete(self.workers, worker.identity)
}

// Decompose
func (self *loadBalancer) decomposeMapOfInstancesMsg(msg []byte) []byte {
	log.Info("----------------------------- Entra en decomposeMapOfInstancesMsg")
	var levels []interface{}
	if err := json.Unmarshal(msg, &levels); err != nil {
		// TODO: Send an error
	}
	var computedInstances []interface{}
	var processMapLevels func([]interface{})
	processMapLevels = func(levels []interface{}) {
		for _, level := range levels {
			if level != nil {
				kind := reflect.TypeOf(level).Kind()
				if kind == reflect.Slice || kind == reflect.Array {
					processMapLevels(level.([]interface{}))
				} else {
					computedInstances = append(computedInstances, levels...)
					return
				}
			}
		}
	}
	processMapLevels(levels)
	// TODO: extract uris
	sortedInstanceUris := make([]string, 0)
	for _, instance := range computedInstances {
		if instance != nil {
			if uri, ok := instance.(map[string]interface{})["Info"].(map[string]interface{})["uri"]; ok {
				sortedInstanceUris = append(sortedInstanceUris, uri.(string))
			}
		} else {
			log.Warn("Remove instance without uri attribute from last balancer response")
		}
	}

	uris, err := json.Marshal(sortedInstanceUris)
	if err != nil {
		// TODO
	}
	return uris
}

// Dispatch chains advancing a shackle
func (self *loadBalancer) advanceShackle(chain lbChain) {
	log.Info("----------------------------- Entra en advanceShackle")
	elem := chain.shackles.Pop()
	// log.Infof("Elem %#v", elem)
	if elem == nil {
		// Decompose
		instanceUrisMsg := self.decomposeMapOfInstancesMsg(chain.msg)
		msg := [][]byte{[]byte(chain.client), nil, instanceUrisMsg}
		// msg := [][]byte{[]byte(chain.client), nil, chain.msg}
		self.frontend.SendMultipart(msg, 0)
		return
	}
	shackle, _ := elem.Value.(lbShackle)
	args, _ := json.Marshal(shackle.serviceArgs)
	msg := [][]byte{[]byte(chain.client), nil, chain.msg, args}
	self.dispatch(self.requireService(shackle.serviceName), msg)
}

// Dispatch requests to waiting workers as possible
func (self *loadBalancer) dispatch(service *lbService, msg [][]byte) {
	log.Info("----------------------------- Entra en dispatch")
	if service == nil {
		log.Fatal("Nil service")
	}
	// Queue message if any
	if len(msg) != 0 {
		service.requests = append(service.requests, msg)
	}
	self.purgeWorkers()
	for service.waiting.Len() > 0 && len(service.requests) > 0 {
		msg, service.requests = service.requests[0], service.requests[1:]
		elem := service.waiting.Pop()
		self.waiting.Delete(elem.Value)
		worker, _ := elem.Value.(*lbWorker)

		self.sendToWorker(worker, SIGNAL_REQUEST, nil, msg)
	}
}

// Register chain from new client request
func (self *loadBalancer) registerChain(client []byte, msg [][]byte) {
	log.Info("----------------------------- Entra en registerChain")
	var services []Balancer
	if err := json.Unmarshal(msg[1], &services); err != nil {
		panic(err)
	}
	chain := lbChain{
		app:      string(msg[0]),
		client:   string(client),
		msg:      msg[2],
		shackles: NewList(),
	}
	for _, service := range services {
		args := service.Args
		args["appId"] = chain.app
		chain.shackles.PushBack(lbShackle{
			serviceName: service.Id,
			// serviceArgs: service.Args,
			serviceArgs: args,
		})
	}
	self.chains[string(client)] = chain
}

// Process a request coming from a client.
func (self *loadBalancer) processClient(client []byte, msg [][]byte) {
	log.Info("----------------------------- Entra en processClient")
	// Application + Services + Instances
	if len(msg) < 3 {
		log.Fatal("Invalid message from client sender")
	}
	// Register chain
	self.registerChain(client, msg)
	// Start chain of requests
	self.advanceShackle(self.chains[string(client)])
}

// Process message sent to us by a worker.
func (self *loadBalancer) processWorker(sender []byte, msg [][]byte) {
	log.Info("----------------------------- Entra en processWorker")
	//  At least, command
	if len(msg) < 1 {
		log.Warn("Invalid message from Worker, this doesn contain command")
		return
	}

	command, msg := msg[0], msg[1:]
	identity := hex.EncodeToString(sender)
	worker, workerReady := self.workers[identity]
	if !workerReady {
		worker = &lbWorker{
			identity: identity,
			address:  sender,
			expiry:   time.Now().Add(HEARTBEAT_EXPIRY),
		}
		self.workers[identity] = worker
	}

	switch string(command) {
	case SIGNAL_READY:
		log.Info("COMMAND: READY")
		//  At least, a service name
		if len(msg) < 1 {
			log.Warn("Invalid message from worker, service name is missing")
			self.deleteWorker(worker, true)
			return
		}
		service := msg[0]
		//  Not first command in session or Reserved service name
		if workerReady || string(service[:4]) == INTERNAL_SERVICE_PREFIX {
			self.deleteWorker(worker, true)
		} else {
			//  Attach worker to service and mark as idle
			worker.service = self.requireService(string(service))
			self.workerWaiting(worker)
		}
	case SIGNAL_REPLY:
		log.Info("COMMAND: REPLY")
		if workerReady {
			//  Remove & save client return envelope and insert the
			//  protocol header and service name, then rewrap envelope.
			client := msg[0]
			chain := self.chains[string(client)]
			chain.msg = msg[2]
			self.advanceShackle(chain)
			self.workerWaiting(worker)
		} else {
			self.deleteWorker(worker, true)
		}
	case SIGNAL_HEARTBEAT:
		log.Info("COMMAND: HEARTBEAT")
		if workerReady {
			worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
		} else {
			self.deleteWorker(worker, true)
		}
	case SIGNAL_DISCONNECT:
		log.Info("COMMAND: DISCONNECT")
		self.deleteWorker(worker, false)
	default:
		log.Info("COMMAND: OTHER")
		log.Warn("Invalid message in Load Balancer")
		Dump(msg)
	}
	return
}

//  Look for & kill expired workers.
//  Workers are oldest to most recent, so we stop at the first alive worker.
func (self *loadBalancer) purgeWorkers() {
	log.Info("----------------------------- Entra en purgeWorkers")
	now := time.Now()
	for elem := self.waiting.Front(); elem != nil; elem = self.waiting.Front() {
		worker, _ := elem.Value.(*lbWorker)
		if worker.expiry.After(now) {
			// TODO: continue
			break
		}
		self.deleteWorker(worker, false)
	}
}

// Locates the service (creates if necessary).
func (self *loadBalancer) requireService(name string) *lbService {
	log.Info("----------------------------- Entra en requireService: " + name)
	if len(name) == 0 {
		log.Warn("Invalid service name have been required")
	}
	// log.Infof("***************************** Service Len: %d", len(self.services))
	// for key, _ := range self.services {
	// 	log.Info(">>>>>>>>>>>>>>>> Service KEY: " + key)
	// }
	service, ok := self.services[name]
	if !ok {
		service = &lbService{
			name:    name,
			waiting: NewList(),
		}
		self.services[name] = service
	}
	return service
}

//  Send message to worker.
//  If message is provided, sends that message.
func (self *loadBalancer) sendToWorker(worker *lbWorker, command string, option []byte, msg [][]byte) {
	log.Info("----------------------------- Entra en sendToWorker")

	//  Stack routing and protocol envelopes to start of message and routing envelope
	if len(option) > 0 {
		msg = append([][]byte{option}, msg...)
	}
	msg = append([][]byte{worker.address, nil, []byte(command)}, msg...)

	// if self.verbose {
	// 	log.Printf("I: sending %X to worker\n", command)
	// 	Dump(msg)
	// }
	self.backend.SendMultipart(msg, 0)
}

//  Handle internal service according to 8/MMI specification
// func (self *server) serviceInternal(service []byte, msg [][]byte) {
// 	// TODO: Change errors code to http erros
// 	returncode := "501"
// 	if string(service) == "mmi.service" {
// 		name := string(msg[len(msg)-1])
// 		if _, ok := self.services[name]; ok {
// 			returncode = "200"
// 		} else {
// 			returncode = "404"
// 		}
// 	}
// 	msg[len(msg)-1] = []byte(returncode)
// 	//  insert the protocol header and service name after the routing envelope
// 	msg = append([][]byte{msg[0], nil, []byte(MDPC_CLIENT), service}, msg[2:]...)
// 	self.socket.SendMultipart(msg, 0)
// }

// This worker is now waiting for work.
func (self *loadBalancer) workerWaiting(worker *lbWorker) {
	//  Queue to broker and service waiting lists
	log.Info("----------------------------- Entra en workerWaiting")
	self.waiting.PushBack(worker)
	worker.service.waiting.PushBack(worker)
	worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
	self.dispatch(worker.service, nil)
}

func (self *loadBalancer) Close() {
	log.Info("----------------------------- Entra en Close")
	if self.frontend != nil {
		self.frontend.Close()
	}
	if self.backend != nil {
		self.backend.Close()
	}
	self.context.Close()
}

// Main broker working loop
func (self *loadBalancer) Run() {
	for {
		items := zmq.PollItems{
			zmq.PollItem{Socket: self.frontend, Events: zmq.POLLIN},
			zmq.PollItem{Socket: self.backend, Events: zmq.POLLIN},
		}

		_, err := zmq.Poll(items, HEARTBEAT_INTERVAL)
		if err != nil {
			log.Fatal("Non items for polling")
		}

		if items[0].REvents&zmq.POLLIN != 0 {
			msg, _ := self.frontend.RecvMultipart(0)
			// TODO: check msg parts
			requestId := msg[0]
			msg = msg[2:]
			self.processClient(requestId, msg)
		}
		if items[1].REvents&zmq.POLLIN != 0 {
			msg, _ := self.backend.RecvMultipart(0)
			sender := msg[0]
			msg = msg[2:]
			self.processWorker(sender, msg)
		}

		if self.heartbeatAt.Before(time.Now()) {
			self.purgeWorkers()
			for elem := self.waiting.Front(); elem != nil; elem = elem.Next() {
				worker, _ := elem.Value.(*lbWorker)
				self.sendToWorker(worker, SIGNAL_HEARTBEAT, nil, nil)
			}
			self.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
		}
	}
}

// func main() {
// 	verbose := len(os.Args) >= 2 && os.Args[1] == "-v"
// 	broker := NewBroker("tcp://*:5555", verbose)
// 	defer broker.Close()

// 	broker.Run()
// }
