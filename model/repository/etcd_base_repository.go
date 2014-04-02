package repository

import (
	"github.com/innotech/hydra/database/connector"
	"github.com/innotech/hydra/model/entity"
	"log"
)

const KEY_PREFIX string = "/db"

type EtcdAccessLayer interface {
	Delete(key string) error
	Get(key string) (*entity.EtcdBaseModel, error)
	GetAll() (*entity.EtcdBaseModels, error)
	GetCollection() string
	Set(entity *entity.EtcdBaseModel) error
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
	prefix := KEY_PREFIX + e.collection
	if key != "" {
		return prefix + "/" + key
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

func (e *EtcdBaseRepository) Set(entity *entity.EtcdBaseModel) error {
	ops, err := entity.ExportEtcdOperations()
	if err != nil {
		log.Fatal("Error expoting etcd operations")
		return err
	}
	// var i = 0
	for key, value := range ops {
		if err := e.conn.Set(e.makePath(key), false, value, PERMANENT); err != nil {
			// TODO: logger
			return err
		}
	}
	return nil
}
