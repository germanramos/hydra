package server

import (
	"net"

	"github.com/innotech/hydra/server/controller"

	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type PublicServer struct {
	controllers                 []controller.Controller
	Listener                    net.Listener
	loadBalancerFrontendAddress string
	Router                      *mux.Router
}

func NewPublicServer(l net.Listener, loadBalancerFrontendAddress string) *PublicServer {
	var p = new(PublicServer)
	p.Listener = l
	p.loadBalancerFrontendAddress = loadBalancerFrontendAddress
	p.Router = mux.NewRouter()
	p.registerControllers()

	return p
}

func (p *PublicServer) registerControllers() {
	p.controllers = make([]controller.Controller, 1)
	p.controllers[0], _ = controller.NewBalancedInstancesController(p.loadBalancerFrontendAddress)
}

func (p *PublicServer) RegisterHandlers() {
	for _, c := range p.controllers {
		c.RegisterHandlers(p.Router)
	}
}
