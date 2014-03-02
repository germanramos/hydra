package server

import (
	"net"

	"github.com/innotech/hydra/server/controller"

	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type PrivateServer struct {
	controllers []controller.Controller
	Listener    net.Listener
	Router      *mux.Router
}

func NewPrivateServer(l net.Listener) *PrivateServer {
	var p = new(PrivateServer)
	p.controllers = make([]controller.Controller, 1)
	p.controllers[0] = controller.NewApplicationController()
	// p.controllers = []*controller.Controller{}
	// p.controllers = [...]controller.Controller{
	// 	controller.NewApplicationController(),
	// }
	p.Listener = l
	p.Router = mux.NewRouter()
	return p
}

func (p *PrivateServer) RegisterControllers() {
	for _, c := range p.controllers {
		c.RegisterHandlers(p.Router)
	}
}

// func (p PrivateServer) Start() {
// 	http.Serve(p.listener, p.router)
// }
