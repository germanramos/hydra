package controller

import (
	"encoding/json"
	"net/http"

	. "github.com/innotech/hydra/load_balancer"
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/model/repository"
	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

// var ZMQ_EMPTY_PART = []byte("")

type BalancedInstancesController struct {
	*BasicController
	loadBalancerAddress string
}

func NewBalancedInstancesController(loadBalancerAddresss string) (*BalancedInstancesController, error) {
	var b = new(BalancedInstancesController)
	b.basePath = "/apps"
	b.loadBalancerAddress = loadBalancerAddresss
	var err error
	b.PathVariables, err = extractPathVariables(b.basePath)
	if err != nil {
		return nil, err
	}
	b.repo = NewEctdRepository()
	b.repo.SetCollection("/apps")
	return b, nil
}

func (b *BalancedInstancesController) sendZMQRequestToBalancer(app []byte, data [][]byte) (reply [][]byte) {
	client := NewClient(b.loadBalancerAddress)
	defer client.Close()

	reply = client.Send(app, data)
	return
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
	balancers, err := json.Marshal(appEntity.Balancers)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	instances, err := json.Marshal(appEntity.Instances)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Request to balancer server
	response := b.sendZMQRequestToBalancer([]byte(appEntity.Id), [][]byte{balancers, instances})
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
