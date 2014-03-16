package balancer

import (
	"errors"
	"http"

	zmq "github.com/alecthomas/gozmq"

	"github.com/innotech/hydra/log"
	"github.com/innotech/hydra/model/entity"
	"github.com/innotech/hydra/model/repository"
)

type BalancerServer struct {
	addr          string
	appRepository *repository.EtcdBaseRepository
	pipelines     map[string]string
}

func NewBalancerServer() {
	b := new(BalancerServer)
	b.appRepository = repository.NewEctdRepository()
	// applicaions global constant or in database module?
	b.appRepository.SetCollection("/applications")
}

func (b *BalancerServer) ExtractAppIDFromClientMultipartMsg(parts [][]byte) (string, error) {
	// TODO 2 to constant
	if len(parts) < 2 {
		return "", errors.New("Multipart message from Client API doesn't contain application ID")
	}
	if len(parts[1]) == 0 {
		return "", errors.New("Multipart message from Client API contains an empty application ID")
	}
	return parts[1], nil
}

func (b BalancerServer) ExtractBalancerPipelineFromApplicationData(data map[string]interface{}) ([]string, error) {

}

func (b *BalancerServer) RegisterPipeline(app *entity.EtcdBaseModel) {
	id, data := app.Explode()
	balancerPipeline := b.ExtractBalancerPipelineFromApplicationData()
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
			appID, err := b.ExtractAppIDFromClientMultipartMsg(parts)
			if err != nil {
				log.Warnf("Unable to balance client request: %s", err)
				// http.Error(rw, err.Error(), http.StatusInternalServerError)
				// response msg contains status request
				// TODO: an error response
			} else {
				app, err := b.appRepository.Get(appID)
				if err != nil {
					log.Warnf("Requested application from hydra client have not been found in database: %s", err)
					// http.Error(rw, err.Error(), http.StatusNotFound)
					// TODO: an error response
					// return
				} else {

					backend.SendMultipart(parts, 0)
				}

			}
		}
		if items[1].REvents&zmq.POLLIN != 0 {
			// Receiving request from balancers
			parts, _ := backend.RecvMultipart(0)

			frontend.SendMultipart(parts, 0)
		}
	}
}
