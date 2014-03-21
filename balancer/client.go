package balancer

import (
	zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"
	uuid "github.com/innotech/hydra/vendors/github.com/nu7hatch/gouuid"
)

// TODO: change to Sender
type Requester interface {
	Close()
	Send([]byte, [][]byte) [][]byte
}

type client struct {
	socket  *zmq.Socket
	context *zmq.Context
	server  string
	timeout time.Duration
}

func NewClient(server string) Client {
	context, _ := zmq.NewContext()
	self := &Client{
		server:  server,
		context: context,
		timeout: 2500 * time.Millisecond,
	}
	self.connect()
	return self
}

func (self *client) connect() {
	if self.socket != nil {
		self.socket.Close()
	}

	self.socket, _ = self.context.NewSocket(zmq.REQ)
	// TODO: I think that uuid is not necessary for Router
	identityUUID, _ := uuid.NewV4()
	identity := identityUUID.String()
	self.socket.SetIdentity(identity)
	self.socket.SetLinger(0)
	self.socket.Connect(self.server)
}

func (self *client) Close() {
	if self.socket != nil {
		self.socket.Close()
	}
	self.context.Close()
}

func (self *client) Send(service []byte, request [][]byte) (reply [][]byte) {
	frame := append([][]byte{service}, request...)

	self.client.SendMultipart(frame, 0)
	msg, _ := self.client.RecvMultipart(0)
	reply = msg[2:]
	return
}
