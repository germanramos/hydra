package controller

import (
	"encoding/json"
	"log"
	"net/http"

	. "github.com/innotech/hydra/load_balancer"
	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/model/repository"
	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

// var ZMQ_EMPTY_PART = []byte("")

type BalancedInstancesController struct {
	BasicController
	loadBalancerAddress string
	requestTimeout      int
}

func NewBalancedInstancesController(loadBalancerAddresss string, requestTimeout int) (*BalancedInstancesController, error) {
	var b = new(BalancedInstancesController)
	b.basePath = "/apps"
	b.loadBalancerAddress = loadBalancerAddresss
	b.requestTimeout = requestTimeout
	var err error
	b.PathVariables, err = extractPathVariables(b.basePath)
	if err != nil {
		return nil, err
	}
	b.repo = NewEctdRepository()
	b.repo.SetCollection(b.basePath)
	return b, nil
}

func (b *BalancedInstancesController) sendZMQRequestToBalancer(app []byte, data [][]byte) (reply [][]byte) {
	// log.Println(b.loadBalancerAddress)
	client := NewClient(b.loadBalancerAddress, b.requestTimeout)
	defer client.Close()

	log.Println("App: " + string(app))
	// Dump(data)
	log.Println("¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨ START sendZMQRequestToBalancer")
	reply = client.Send(app, data)
	log.Println("¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨ END sendZMQRequestToBalancer")
	// Dump(reply)
	// log.Println("¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨ END RESPONSE sendZMQRequestToBalancer")
	return
}

func (b *BalancedInstancesController) RegisterHandlers(r *mux.Router) {
	// r.HandleFunc(b.basePath, b.List).Methods("GET")
	r.HandleFunc(b.basePath+"/{id}", b.Get).Methods("GET")
	// retro compatibility alias
	r.HandleFunc("/app/{id}", b.Get).Methods("GET")
}

func (b *BalancedInstancesController) getActiveInstances(instances []Instance) []Instance {
	activeInstances := make([]Instance, 0)
	for _, instance := range instances {
		if len(instance.Info) > 0 && instance.Info["state"] == "0.00" {
			activeInstances = append(activeInstances, instance)
		}
	}
	return activeInstances
}

func (b *BalancedInstancesController) Get(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appId := vars["id"]
	app, err := b.repo.Get(appId)
	// log.Printf("Repo Get Request %#v", map[string]interface{}(*app))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	appEntity, err := NewApplicationFromEtcdBaseModel(app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// log.Printf("appEntity: %#v", appEntity)
	// log.Printf("appEntity.Balancers: %#v", appEntity.Balancers)

	var jsonOutput []byte = []byte("[]")
	if len(appEntity.Balancers) > 0 {
		balancers, err := json.Marshal(appEntity.Balancers)
		if err != nil {
			log.Println("Bad format for Balancers")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		activeInstances := b.getActiveInstances(appEntity.Instances)
		if len(activeInstances) > 0 {
			instances, err := json.Marshal(activeInstances)
			if err != nil {
				log.Println("Bad format for Instances")
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			response := b.sendZMQRequestToBalancer([]byte(appEntity.Id), [][]byte{balancers, instances})
			log.Println("++++++++++++++++++ ZMQ RESPONSE ++++++++++++++++++")
			// log.Printf("%#v", response)
			// TODO: process response

			if len(response) > 0 {
				jsonOutput = response[0]
			} else {
				log.Println("Zeromq request timeout expired")
				http.Error(rw, "Zeromq request timeout expired", http.StatusInternalServerError)
				return
			}

		} else {
			jsonOutput = []byte("[]")
		}
	} else {
		activeInstances := b.getActiveInstances(appEntity.Instances)
		sortedInstanceUris := make([]string, 0)
		for _, instance := range activeInstances {
			sortedInstanceUris = append(sortedInstanceUris, instance.Info["uri"].(string))
		}

		jsonOutput, _ = json.Marshal(sortedInstanceUris)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}
