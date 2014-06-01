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
	b.repo.SetCollection(b.basePath)
	return b, nil
}

func (b *BalancedInstancesController) sendZMQRequestToBalancer(app []byte, data [][]byte) (reply [][]byte) {
	// log.Println(b.loadBalancerAddress)
	client := NewClient(b.loadBalancerAddress)
	defer client.Close()

	log.Println("App: " + string(app))
	// Dump(data)
	reply = client.Send(app, data)
	// log.Println("¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨¨ RESPONSE sendZMQRequestToBalancer")
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
		if len(instance.Info) > 0 {
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
			// TODO: process response

			jsonOutput = response[0]
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

	// balancers, err := json.Marshal(appEntity.Balancers)
	// // log.Printf("Balancers: " + string(balancers))
	// if err != nil {
	// 	log.Println("No Balancers")
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // TODO: check si hay Balanceadores a los que enviar por zeromq

	// instances, err := json.Marshal(appEntity.Instances)
	// // log.Printf("Instances: " + string(instances))
	// if err != nil {
	// 	log.Println("No Instances")
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // TODO: check si hay instancias y si hay eliminar las vacias antes de enviar al load balancer

	// // log.Println("PRE ZMQ send request")
	// // log.Println(appId)
	// // log.Println(appEntity.Id)
	// // Request to balancer server
	// response := b.sendZMQRequestToBalancer([]byte(appEntity.Id), [][]byte{balancers, instances})
	// // TODO: process response

	// jsonOutput := response[0]
	// // jsonOutput, err := json.Marshal(response)
	// // log.Println("EMIT RESPONSE TO FINAL CLIENT: " + string(jsonOutput))
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}
