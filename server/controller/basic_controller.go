package controller

import (
	"encoding/json"
	"errors"
	"github.com/innotech/hydra/log"
	// "fmt"
	"net/http"
	// "strconv"
	"strings"

	"github.com/innotech/hydra/model/entity"
	"github.com/innotech/hydra/model/repository"

	"github.com/innotech/hydra/vendors/github.com/gorilla/mux"
)

type Controller interface {
	GetConfiguredRepository(pathVars map[string]string) *repository.EtcdBaseRepository
	RegisterHandlers(r *mux.Router)
}

type BasicController struct {
	basePath      string
	PathVariables []string
	repo          *repository.EtcdBaseRepository
	setValidation func(map[string]interface{}, map[string]string) bool
}

func NewBasicController(basePath string, setValidation func(map[string]interface{}, map[string]string) bool) (*BasicController, error) {
	var b = new(BasicController)
	b.basePath = basePath
	b.setValidation = setValidation
	var err error
	b.PathVariables, err = extractPathVariables(basePath)
	if err != nil {
		return nil, err
	}
	b.repo = repository.NewEctdRepository()
	return b, nil
}

func getBoundariesIndexesForNextPathVariable(path string) (i1, i2 int) {
	i1 = strings.Index(path, "{")
	i2 = strings.Index(path, "}")
	return
}

func extractPathVariables(path string) ([]string, error) {
	var variables []string
	i1, i2 := getBoundariesIndexesForNextPathVariable(path)
	for i1 != -1 || i2 != -1 {
		if (i1 != -1 && i2 == -1) || (i1 == -1 && i2 != -1) {
			return nil, errors.New("Invalid controller path: ill-defined variables")
		}
		variables = append(variables, path[i1+1:i2])
		path = path[i2+1:]
		i1, i2 = getBoundariesIndexesForNextPathVariable(path)
	}
	return variables, nil
}

func (b *BasicController) GetConfiguredRepository(pathVars map[string]string) *repository.EtcdBaseRepository {
	finalPath := b.basePath
	for key, value := range pathVars {
		finalPath = strings.Replace(finalPath, "{"+key+"}", value, 1)
	}
	log.Info("Controller Set Collection " + finalPath)
	b.repo.SetCollection(finalPath)
	return b.repo
}

func (a *BasicController) RegisterHandlers(r *mux.Router) {
	r.HandleFunc(a.basePath+"/{id}", a.Delete).Methods("DELETE")
	r.HandleFunc(a.basePath+"/{id}", a.Get).Methods("GET")
	r.HandleFunc(a.basePath, a.List).Methods("GET")
	r.HandleFunc(a.basePath, a.Set).Methods("POST")
}

func (a *BasicController) Delete(rw http.ResponseWriter, req *http.Request) {
	// TODO
}

func (a *BasicController) Get(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	app, err := a.GetConfiguredRepository(vars).Get(id)
	// TODO: Implement Balancer Middleware
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	jsonOutput, err := json.Marshal(app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}

func (a *BasicController) List(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	apps, err := a.GetConfiguredRepository(vars).GetAll()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	jsonOutput, err := json.Marshal(apps)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonOutput)
}

func (a *BasicController) Set(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var app entity.EtcdBaseModel
	err := decoder.Decode(&app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	vars := mux.Vars(req)

	// log.Infof("%#v", app)
	// log.Infof("%#v", vars)
	if a.setValidation(app, vars) != true {
		log.Warn("Post Instance: Set validation fail")
		http.Error(rw, "Post Instance: Set validation fail", http.StatusBadRequest)
		return
	}

	err = a.GetConfiguredRepository(vars).Set(&app)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
