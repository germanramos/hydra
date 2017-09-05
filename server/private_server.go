package server

import (
	"net"

	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/model/repository"
	"github.com/innotech/hydra/server/controller"
	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type PrivateServer struct {
	controllers []controller.Controller
	Listener    net.Listener
	Router      *mux.Router
}

func NewPrivateServer(l net.Listener, defaultTTL int) *PrivateServer {
	var p = new(PrivateServer)
	p.Listener = l
	p.Router = mux.NewRouter()
	p.registerControllers(defaultTTL)

	return p
}

func validateApp(app map[string]interface{}, vars map[string]string) bool {
	return true
}

func validateInstance(app map[string]interface{}, vars map[string]string) bool {
	if len(app) == 1 {
		var repo *repository.EtcdBaseRepository = repository.NewEctdRepository()
		repo.SetCollection("/apps")
		_, err := repo.Get(vars["appId"])
		if err != nil {
			log.Warn(err)
			log.Warn("validateInstance: Error getting app " + vars["appId"])
			return false
		}
		return true
	} else {
		return false
	}
}

func (p *PrivateServer) registerControllers(defaultTTL int) {
	p.controllers = make([]controller.Controller, 2)
	p.controllers[0], _ = controller.NewBasicController("/apps", true, defaultTTL, validateApp)
	p.controllers[1], _ = controller.NewBasicController("/apps/{appId}/Instances", true, defaultTTL, validateInstance)
}

func (p *PrivateServer) RegisterHandlers() {
	for _, c := range p.controllers {
		c.RegisterHandlers(p.Router)
	}
}
