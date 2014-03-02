package controller

import (
	"encoding/json"
	"net/http"

	"github.com/innotech/hydra/model/entity"
	"github.com/innotech/hydra/model/repository"

	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type ApplicationController struct {
	basePath string
	repo     *repository.EtcdBaseRepository
}

func NewApplicationController() *ApplicationController {
	var a = new(ApplicationController)
	a.basePath = "/applications/"
	a.repo = repository.NewEctdRepository()
	return a
}

func (a ApplicationController) RegisterHandlers(r *mux.Router) {
	r.HandleFunc(a.basePath+"{id}", a.Delete).Methods("DELETE")
	r.HandleFunc(a.basePath+"{id}", a.Get).Methods("GET")
	r.HandleFunc(a.basePath, a.List).Methods("GET")
	r.HandleFunc(a.basePath, a.Set).Methods("POST")
}

func (a ApplicationController) Delete(rw http.ResponseWriter, req *http.Request) {

}

func (a ApplicationController) Get(rw http.ResponseWriter, req *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

}

func (a ApplicationController) List(rw http.ResponseWriter, req *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

}

func (a *ApplicationController) Set(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var app entity.EtcdBaseModel
	err := decoder.Decode(&app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.repo.Set(&app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
