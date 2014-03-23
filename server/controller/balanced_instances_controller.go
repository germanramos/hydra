package controller

import (
	// "encoding/json"
	// "net/http"

	// zmq "github.com/innotech/hydra/vendors/github.com/alecthomas/gozmq"
	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
	// uuid "github.com/innotech/hydra/vendors/github.com/nu7hatch/gouuid"

	// "github.com/innotech/hydra/load_balancer"
	// . "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/model/repository"
)

var ZMQ_EMPTY_PART = []byte("")

type BalancedInstancesController struct {
	*BasicController
}

func NewBalancedInstancesController() (*BalancedInstancesController, error) {
	var b = new(BalancedInstancesController)
	b.basePath = "/applications"
	var err error
	b.PathVariables, err = extractPathVariables(b.basePath)
	if err != nil {
		return nil, err
	}
	b.repo = NewEctdRepository()
	// TODO: fixed routes to constants
	b.repo.SetCollection("/applications")
	return b, nil
}

// func (b *BalancedInstancesController) sendZMQRequestToBalancer(app []byte, data [][]byte) (reply [][]byte) {
// 	// TODO: Load Balancer from config
// 	client := NewClient("tcp://localhost:5555" /*, verbose*/)
// 	defer client.Close()

// 	reply = client.Send(app, data)
// 	return

// 	// context, _ := zmq.NewContext()
// 	// defer context.Close()

// 	// // Set unique identity to make tracing possible
// 	// identityUUID, _ := uuid.NewV4()
// 	// identity := identityUUID.String()

// 	// client, _ := context.NewSocket(zmq.REQ)
// 	// client.SetIdentity(identity)
// 	// // TODO: Make a constant address
// 	// client.Connect("ipc://frontend.ipc")
// 	// defer client.Close()

// 	// // Send request, get reply
// 	// client.SendMultipart([][]byte{[]byte(identity), ZMQ_EMPTY_PART, jsonApp}, 0)
// 	// reply, _ := client.Recv(0)
// 	// return reply
// }

func (b *BalancedInstancesController) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/applications/{appId}/instances", b.List).Methods("GET")
}

// func (b *BalancedInstancesController) List(rw http.ResponseWriter, req *http.Request) {
// 	vars := mux.Vars(req)
// 	appId := vars["appId"]
// 	app, err := b.repo.Get(appId)
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	appEntity, err := NewApplicationFromEtcdBaseModel(app)
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	balancers, err := json.Marshal(appEntity.Balancers)
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	instances, err := json.Marshal(appEntity.Instances)
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Request to balancer server
// 	// response := b.sendZMQRequestToBalancer([]byte(appEntity.Id), [][]byte{balancers, instances})
// 	// TODO: process response

// 	jsonOutput, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(rw, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	rw.WriteHeader(http.StatusOK)
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.Write(jsonOutput)
// }
