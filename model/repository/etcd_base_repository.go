package repository

import (
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/model/entity"
	"log"
	"net/http"
	// "time"
)

const KEY_PREFIX string = "/db"

type EtcdAccessLayer interface {
	Delete(key string) error
	Get(key string) (*entity.EtcdBaseModel, error)
	GetAll() (*entity.EtcdBaseModels, error)
	GetCollection() string
	Set(entity *entity.EtcdBaseModel, ttl string, w http.ResponseWriter, req *http.Request) error
	SetCollection(collection string)
}

type EtcdBaseRepository struct {
	conn       *connector.EtcdConnector
	collection string
}

func NewEctdRepository() *EtcdBaseRepository {
	var e = new(EtcdBaseRepository)
	e.conn = connector.GetEtcdConnector()
	return e
}

func (e *EtcdBaseRepository) GetCollection() string {
	return e.collection
}

func (e *EtcdBaseRepository) SetCollection(collection string) {
	e.collection = collection
}

func (e *EtcdBaseRepository) makePath(key string) string {
	// log.Println(">>>>>>>>>>>>>>>> KEY_PREFIX: " + KEY_PREFIX)
	// log.Println(">>>>>>>>>>>>>>>> e.collection: " + e.collection)
	prefix := KEY_PREFIX + e.collection
	if key != "" {
		if string(key[0]) == "/" {
			return prefix + key
		} else {
			return prefix + "/" + key
		}
		// log.Println(">>>>>>>>>>>>>>>> PREFIX: " + prefix)
		// // return prefix + "/" + key
		// log.Println(">>>>>>>>>>>>>>>> FINAL KEY: " + prefix + key)
		// return prefix + key
	}
	return prefix
}

func (e *EtcdBaseRepository) Delete(key string) error {
	// TODO
	return nil
}

func (e *EtcdBaseRepository) Get(key string) (*entity.EtcdBaseModel, error) {
	event, err := e.conn.Get(e.makePath(key), true, false)
	if err != nil {
		return nil, err
	}
	return entity.NewModelFromEvent(event)
}

func (e *EtcdBaseRepository) GetAll() (*entity.EtcdBaseModels, error) {
	event, err := e.conn.Get(e.makePath(""), true, false)
	if err != nil {
		return nil, err
	}
	return entity.NewModelsFromEvent(event)
}

func (e *EtcdBaseRepository) Set(entity *entity.EtcdBaseModel, ttl string, w http.ResponseWriter, req *http.Request) error {
	ops, err := entity.ExportEtcdOperations()
	if err != nil {
		log.Fatal("Error expoting etcd operations")
		return err
	}
	for key, value := range ops {
		log.Println("Probe KEY: " + e.makePath(key))
		log.Println("Probe KEY TTL: " + ttl)
		var dir bool = false
		if value == "" {
			dir = true
		}
		log.Printf("DIR: %#v", dir)
		if err := e.conn.Set(e.makePath(key), dir, value, ttl, w, req); err != nil {
			// TODO: logger
			log.Println("SET ERROR")
			log.Println(err)
			return err
		}
	}
	return nil
}
