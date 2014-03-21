package balancer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"http"

	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"

	"github.com/innotech/hydra/log"
	// . "github.com/innotech/hydra/model/entity"
)

const (
	INTERNAL_SERVICE_PREFIX = "mmi."
	// Merge all heartbeat
	HEARTBEAT_INTERVAL = 2500 * time.Millisecond
	HEARTBEAT_EXPIRY   = HEARTBEAT_INTERVAL * HEARTBEAT_LIVENESS
)

type Broker interface {
	Close()
	Run()
}

type worker struct {
	identity string    //  UUID Identity of worker
	address  []byte    //  Address to route to
	expiry   time.Time //  Expires at this point, unless heartbeat
	service  *service  //  Owning service, if known
}

// shackle

type pipeline struct {
	request string
	workers []string
}

type service struct {
	broker   Broker
	name     string
	requests [][][]byte //  List of client requests
	waiting  *ZList     //  List of waiting workers
}

type loadBalancer struct {
	context     *zmq.Context        //  Context
	heartbeatAt time.Time           //  When to send HEARTBEAT
	services    map[string]*service //  Known services
	frontend    *zmq.Socket         //  Socket for clients
	backend     *zmq.Socket         //  Socket for workers
	waiting     *ZList              //  Idle workers
	workers     map[string]*worker  //  Known workers
}

func NewLoadBalancer(frontendEndpoint, backendEndpoint string) Broker {
	context, _ := zmq.NewContext()
	// Define inproc socket to talk with HTTP CLient API
	frontend, _ := context.NewSocket(zmq.ROUTER)
	frontend.SetLinger(0)
	frontend.Bind(frontendEndpoint)
	// Define tcp socket to talk with Workers
	backend, _ := context.NewSocket(zmq.ROUTER)
	backend.SetLinger(0)
	backend.Bind(backendEndpoint)
	log.Infof("Load Balancer is active at %s\n", backendEndpoint)
	return &server{
		context:     context,
		heartbeatAt: time.Now().Add(HEARTBEAT_INTERVAL),
		services:    make(map[string]*service),
		frontend:    frontend,
		backend:     backend,
		waiting:     NewList(),
		workers:     make(map[string]*worker),
	}
}

// Deletes worker from all data structures, and deletes worker.
func (self *server) deleteWorker(worker *worker, disconnect bool) {
	if worker == nil {
		log.Warn("Nil worker")
	}

	if disconnect {
		self.sendToWorker(worker, SIGNAL_DISCONNECT, nil, nil)
	}

	if worker.service != nil {
		worker.service.waiting.Delete(worker)
	}
	self.waiting.Delete(worker)
	delete(self.workers, worker.identity)
}

// Dispatch requests to waiting workers as possible
func (self *server) dispatch(service *service, msg [][]byte) {
	if service == nil {
		log.Warn("Nil service")
	}
	//  Queue message if any
	if len(msg) != 0 {
		service.requests = append(service.requests, msg)
	}
	self.purgeWorkers()
	for service.waiting.Len() > 0 && len(service.requests) > 0 {
		msg, service.requests = service.requests[0], service.requests[1:]
		elem := service.waiting.Pop()
		self.waiting.Remove(elem)
		worker, _ := elem.Value.(*mdbWorker)
		self.sendToWorker(worker, MDPW_REQUEST, nil, msg)
	}
}

// Process a request coming from a client.
func (self *server) processClient(sender []byte, msg [][]byte) {
	// Balancer + Instances
	if len(msg) < 2 {
		log.Warn("Invalid message from client requester")
	}
	service := msg[0]
	// Set reply return address to client sender
	msg = append([][]byte{sender, nil}, msg[2:]...)
	if string(service[:4]) == INTERNAL_SERVICE_PREFIX {
		self.serviceInternal(service, msg)
	} else {
		self.dispatch(self.requireService(string(service)), msg)
	}
}

// Process message sent to us by a worker.
func (self *loadBalancer) processWorker(sender []byte, msg [][]byte) {
	//  At least, command
	if len(msg) < 1 {
		log.Warn("Invalid message from Worker")
	}

	command, msg := msg[0], msg[1:]
	// TODO: Why?
	identity := hex.EncodeToString(sender)
	worker, workerReady := self.workers[identity]
	if !workerReady {
		worker = &worker{
			// TODO: Why identity and address
			identity: identity,
			address:  sender,
			expiry:   time.Now().Add(HEARTBEAT_EXPIRY),
		}
		self.workers[identity] = worker
		if self.verbose {
			log.Infof("Registering new worker: %s\n", identity)
		}
	}

	switch string(command) {
	case SIGNAL_READY:
		//  At least, a service name
		if len(msg) < 1 {
			log.Warn("Invalid message from worker, service name is missing")
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
		if workerReady {
			//  Remove & save client return envelope and insert the
			//  protocol header and service name, then rewrap envelope.
			client := msg[0]
			msg = append([][]byte{client, nil, []byte(MDPC_CLIENT), []byte(worker.service.name)}, msg[2:]...)
			self.socket.SendMultipart(msg, 0)
			self.workerWaiting(worker)
		} else {
			self.deleteWorker(worker, true)
		}
	case SIGNAL_HEARTBEAT:
		if workerReady {
			worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
		} else {
			self.deleteWorker(worker, true)
		}
	case SIGNAL_DISCONNECT:
		self.deleteWorker(worker, false)
	default:
		log.Println("E: invalid message:")
		Dump(msg)
	}
}

//  Look for & kill expired workers.
//  Workers are oldest to most recent, so we stop at the first alive worker.
func (self *server) purgeWorkers() {
	now := time.Now()
	for elem := self.waiting.Front(); elem != nil; elem = self.waiting.Front() {
		worker, _ := elem.Value.(*worker)
		if worker.expiry.After(now) {
			// TODO: continue
			break
		}
		self.deleteWorker(worker, false)
	}
}

//  Locates the service (creates if necessary).
func (self *server) requireService(name string) *service {
	if len(name) == 0 {
		log.Warn("Invalid service name have been required")
	}
	service, ok := self.services[name]
	if !ok {
		service = &service{
			name:    name,
			waiting: NewList(),
		}
		self.services[name] = service
	}
	return service
}

//  Send message to worker.
//  If message is provided, sends that message.
func (self *server) sendToWorker(worker *worker, command string, option []byte, msg [][]byte) {
	//  Stack routing and protocol envelopes to start of message and routing envelope
	if len(option) > 0 {
		msg = append([][]byte{option}, msg...)
	}
	msg = append([][]byte{worker.address, nil, []byte(MDPW_WORKER), []byte(command)}, msg...)

	// if self.verbose {
	// 	log.Printf("I: sending %X to worker\n", command)
	// 	Dump(msg)
	// }
	self.socket.SendMultipart(msg, 0)
}

//  Handle internal service according to 8/MMI specification
func (self *server) serviceInternal(service []byte, msg [][]byte) {
	// TODO: Change errors code to http erros
	returncode := "501"
	if string(service) == "mmi.service" {
		name := string(msg[len(msg)-1])
		if _, ok := self.services[name]; ok {
			returncode = "200"
		} else {
			returncode = "404"
		}
	}
	msg[len(msg)-1] = []byte(returncode)
	//  insert the protocol header and service name after the routing envelope
	msg = append([][]byte{msg[0], nil, []byte(MDPC_CLIENT), service}, msg[2:]...)
	self.socket.SendMultipart(msg, 0)
}

// This worker is now waiting for work.
func (self *server) workerWaiting(worker *worker) {
	//  Queue to broker and service waiting lists
	self.waiting.PushBack(worker)
	worker.service.waiting.PushBack(worker)
	worker.expiry = time.Now().Add(HEARTBEAT_EXPIRY)
	self.dispatch(worker.service, nil)
}

func (self *server) Close() {
	if self.socket != nil {
		self.socket.Close()
	}
	self.context.Close()
}

// Main broker working loop
func (self *server) Run() {
	for {
		items := zmq.PollItems{
			zmq.PollItem{Socket: self.frontend, Events: zmq.POLLIN},
			zmq.PollItem{Socket: self.backend, Events: zmq.POLLIN},
		}

		_, err := zmq.Poll(items, HEARTBEAT_INTERVAL)
		if err != nil {
			log.Warn("Non items for polling")
		}

		if items[0].REvents&zmq.POLLIN != 0 {
			msg, _ := frontend.RecvMultipart(0)
			// TODO: check msg parts
			sender := msg[0]
			msg = msg[2:]
			self.processClient(sender, msg)
		}
		if items[1].REvents&zmq.POLLIN != 0 {
			msg, _ := backend.RecvMultipart(0)
			sender := msg[0]
			msg = msg[3:]
			self.processWorker(sender, msg)
		}

		// if item := items[0]; item.REvents&zmq.POLLIN != 0 {
		// 	msg, _ := self.socket.RecvMultipart(0)

		// 	sender := msg[0]
		// 	header := msg[2]
		// 	msg = msg[3:]

		// 	if string(header) == MDPC_CLIENT {
		// 		self.processClient(sender, msg)
		// 	} else if string(header) == MDPW_WORKER {
		// 		self.processWorker(sender, msg)
		// 	} else {
		// 		log.Warn("Load Balancer receive invalid message")
		// 	}
		// }

		if self.heartbeatAt.Before(time.Now()) {
			self.purgeWorkers()
			for elem := self.waiting.Front(); elem != nil; elem = elem.Next() {
				worker, _ := elem.Value.(*worker)
				self.sendToWorker(worker, SIGNAL_HEARTBEAT, nil, nil)
			}
			self.heartbeatAt = time.Now().Add(HEARTBEAT_INTERVAL)
		}
	}
}

func main() {
	verbose := len(os.Args) >= 2 && os.Args[1] == "-v"
	broker := NewBroker("tcp://*:5555", verbose)
	defer broker.Close()

	broker.Run()
}

//////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////

var ZMQ_EMPTY_PART = []byte("")

type BalancerServer struct {
	addr      string
	pipelines map[string][]Balancer
}

func NewBalancerServer() *BalancerServer {
	b := new(BalancerServer)
	return b
}

func (b *BalancerServer) checkPipelineStatus(pipelineId string) bool {

}

func (b *BalancerServer) registerPipeline(requestIdentity string, app *entity.Application) string {
	b.pipelines[requestIdentity] = app.Balancers
	return requestIdentity
}

func (b *BalancerServer) Start() {
	context, _ := zmq.NewContext()
	defer context.Close()

	// Frontend socket talks to client API over inproc
	frontend, _ := context.NewSocket(zmq.ROUTER)
	frontend.Bind("ipc://frontend.ipc")
	defer frontend.Close()

	// Backend socket talks to workers over tcp
	backend, _ := context.NewSocket(zmq.ROUTER)
	backend.Bind("tcp://" + s.addr)
	defer backend.Close()

	// Connect backend to frontend via a proxy
	items := zmq.PollItems{
		zmq.PollItem{Socket: frontend, Events: zmq.POLLIN},
		zmq.PollItem{Socket: backend, Events: zmq.POLLIN},
	}

	for {
		_, err := zmq.Poll(items, -1)
		if err != nil {
			log.Fatalf("Balancer server exited with error:", err)
		}

		if items[0].REvents&zmq.POLLIN != 0 {
			// Receiving request from client API
			parts, _ := frontend.RecvMultipart(0)
			var app entity.Application
			err := json.Unmarshal(parts[3], &app)
			if err != nil {
				// TODO: return error to client
			}
			pipelineId := b.registerPipeline(parts[0], app)
			if ok := checkPipelineStatus(pipelineId); !ok {
				// TODO: Altenative plan
				// TODO: wait or not
				// TODO: send error to client
			}
			// Start pipeline sending data to first balancers
			backend.SendMultipart([][]byte{[]byte(pipelineId), ZMQ_EMPTY_PART, parts[2]}, 0)
		}
		if items[1].REvents&zmq.POLLIN != 0 {
			// Receiving request from balancers
			parts, _ := backend.RecvMultipart(0)
			nextBalancer, thereIsMore := b.pipelines[parts[0]]
			if !thereIsMore {
				frontend.SendMultipart(parts[2], 0)
			} else {
				backend.SendMultipart([][]byte{[]byte(pipelineId), ZMQ_EMPTY_PART, parts[2]}, 0)
			}
		}
	}
}
