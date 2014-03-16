package controller

import (
	"encoding/json"
	"net/http"

	zmq "github.com/alecthomas/gozmq"
	uuid "github.com/nu7hatch/gouuid"

	"github.com/innotech/hydra/model/entity"
)

const ZMQ_EMPTY_PART []byte = []byte("")

type BalancedInstancesController struct {
	*BasicController
}

func NewBalancedInstancesController() (*BasicController, error) {
	var b = new(BalancedInstancesController)
	b.basePath = "/applications"
	var err error
	b.PathVariables, err = extractPathVariables(basePath)
	if err != nil {
		return nil, err
	}
	b.repo = repository.NewEctdRepository()
	// TODO: fixed routes to constants
	b.repo.SetCollection("/applications")
	return b, nil
}

func (b *BalancedInstancesController) sendZMQRequestToBalancer(jsonApp []byte) {
	context, _ := zmq.NewContext()
	defer context.Close()

	// Set unique identity to make tracing possible
	identity, _ := uuid.NewV4()

	client, _ := context.NewSocket(zmq.REQ)
	client.SetIdentity(identity)
	// TODO: Make a constant address
	client.Connect("ipc://frontend.ipc")
	defer client.Close()

	// Send request, get reply
	client.SendMultipart([][]byte{identity, ZMQ_EMPTY_PART, jsonApp}, 0)
	reply, _ := client.Recv(0)
}

func (b *BalancedInstancesController) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/applications/{appId}/instances", b.List).Methods("GET")
}

func (b *BalancedInstancesController) List(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appId := vars["appId"]
	app, err := b.repo.Get(appId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	appEntity, err := NewApplicationFromEtcdBaseModel(app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonApp, err := json.Marshal(appEntity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Request to balancer server
	response := sendZMQRequest(jsonApp)
	// TODO: process response

	jsonOutput, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}
