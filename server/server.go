package server

import (
	// "bytes"
	"encoding/json"
	// "fmt"
	// "io/ioutil"
	"log"
	"net/http"

	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"

	"github.com/innotech/hydra/driver"
	"github.com/innotech/hydra/server/model"
	"github.com/innotech/hydra/server/repository"
)

type server struct {
	instanceRepository *repository.InstanceRepository
	// driver             *driver.EtcdDriver
}

func NewServer(driver *driver.EtcdDriver) *server {
	s := new(server)
	// s.driver = driver
	s.instanceRepository = repository.NewInstaceRepository(driver)
	return s
}

func (s *server) PutInstance(rw http.ResponseWriter, req *http.Request) error {
	// vars := mux.Vars(req)
	// appId := vars["appId"]
	// fmt.Println("appId: " + appId)
	// // var b = []byte(`{"testapp":{"Server":"api"}}`)
	// // var b = []byte(`{
	// // 	"key1": {
	// // 		"Server": "http://mycompany.com/api"
	// // 	}
	// // }`)
	// // var input = bytes.NewReader(b)
	// decoder := json.NewDecoder(r.Body)
	// // body, err := ioutil.ReadAll(r.Body)
	// // if err != nil {
	// // 	fmt.Println("ERROR-----")
	// // 	fmt.Fprintf(w, "%s", err)
	// // }
	// // decoder := json.NewDecoder(input)
	// var instance model.InstanceModel
	// err := decoder.Decode(&instance)
	// // err = json.Unmarshal(body, &instance)
	// fmt.Println(instance["testapp"].Server)
	// if err != nil {
	// 	fmt.Println("ERROR 1 -----")
	// 	log.Println(err)
	// 	fmt.Println("ERROR 1 -----")
	// 	// TODO
	// 	// panic()
	// 	// http.Error(w, "Invalid Instance 1", http.StatusBadRequest)
	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	return nil
	// }

	// // log.Println(t.Test)
	// err = s.instanceRepository.Set(instance, appId)
	// if err != nil {
	// 	// http.Error(w, "Invalid Instance 2", http.StatusBadRequest)
	// 	// w.WriteHeader(http.StatusNoContent)
	// 	return nil
	// }
	// // w.WriteHeader(http.StatusOK)

	// return nil

	decoder := json.NewDecoder(req.Body)
	var i model.InstanceModel
	err := decoder.Decode(&i)
	if err != nil {
		// http.NotFound(rw, req)
		// log.Fatal(err)
		return err
	}
	log.Println(i["testapp"].Server)
	// err = s.instanceRepository.Set(&i, appId)
	// if err != nil {
	// 	// http.Error(w, "Invalid Instance 2", http.StatusBadRequest)
	// 	// w.WriteHeader(http.StatusNoContent)
	// 	return err
	// }
	return nil
}

func (s *server) handleFunc(router *mux.Router, path string, f func(w http.ResponseWriter, r *http.Request) error) *mux.Route {
	return router.HandleFunc(path, func(rw http.ResponseWriter, req *http.Request) {
		// hah, err := ioutil.ReadAll(r.Body)

		// if err != nil {
		// 	fmt.Println("ERROR-----")
		// 	fmt.Fprintf(w, "%s", err)
		// }
		// fmt.Println(string(hah))
		// fmt.Println("HOLAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		// fmt.Fprintf(w, "%s", hah)

		if err := f(rw, req); err != nil {
			// TODO:
			http.NotFound(rw, req)
			log.Fatal(err)
		}

		// decoder := json.NewDecoder(req.Body)
		// var i model.InstanceModel
		// err := decoder.Decode(&i)
		// if err != nil {
		// 	http.NotFound(rw, req)
		// 	log.Fatal(err)
		// }
		// log.Println(i["testapp"].Server)
	})
}

func P(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var i model.InstanceModel
	err := decoder.Decode(&i)
	if err != nil {
		http.NotFound(rw, req)
		log.Fatal(err)
	}
	log.Println(i["testapp"].Server)
}

func (s *server) LoadRouter() /* *mux.Router*/ {
	router := mux.NewRouter()
	s.handleFunc(router, "/applications/{appId}/instances/{serverId}", s.PutInstance).Methods("PUT")
	// router.HandleFunc("/applications/{appId}/instances/{serverId}", P).Methods("PUT")
	http.Handle("/", router)
	// return router
}

func (s *server) Start() {
	// r := s.LoadRouter()
	s.LoadRouter()
	http.ListenAndServe(":8082", nil)
}
