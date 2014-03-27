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
	p.Listener = l
	p.Router = mux.NewRouter()
	p.registerControllers()

	return p
}

func (p *PrivateServer) registerControllers() {
	p.controllers = make([]controller.Controller, 2)
	p.controllers[0], _ = controller.NewBasicController("/apps")
	p.controllers[1], _ = controller.NewBasicController("/apps/{appId}/instances")
}

func (p *PrivateServer) RegisterHandlers() {
	for _, c := range p.controllers {
		c.RegisterHandlers(p.Router)
	}
}
