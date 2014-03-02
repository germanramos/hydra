package repository

import (
	"github.com/innotech/hydra/database/connector"
	// "github.com/innotech/hydra/log"
	"github.com/innotech/hydra/model/entity"
)

type EtcdAccessLayer interface {
	Delete(key string) error
	Get(key string) *entity.EtcdBaseModel
	GetAll() []*entity.EtcdBaseModel
	Set(entity *entity.EtcdBaseModel) error
}

type EtcdBaseRepository struct {
	conn *connector.EtcdConnector
}

func NewEctdRepository() *EtcdBaseRepository {
	var e = new(EtcdBaseRepository)
	e.conn = connector.GetEtcdConnector()
	return e
}

func (e EtcdBaseRepository) Delete(key string) error {
	return nil
}

func (e EtcdBaseRepository) Get(key string) *entity.EtcdBaseModel {
	return nil
}

func (e EtcdBaseRepository) GetAll() []*entity.EtcdBaseModel {
	return nil
}

func (e EtcdBaseRepository) Set(entity *entity.EtcdBaseModel) error {
	ops, err := entity.ExportEtcdOperations()
	if err != nil {
		return err
	}
	// var i = 0
	for key, value := range ops {
		if err := e.conn.Set(key, false, value, PERMANENT); err != nil {
			// TODO: logger
			return err
		}
	}
	return nil
}
