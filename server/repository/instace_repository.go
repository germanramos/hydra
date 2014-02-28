package repository

import (
	"github.com/innotech/hydra/driver"
	"github.com/innotech/hydra/server/model"
	. "github.com/innotech/hydra/utils"

	"time"
)

type InstanceRepository struct {
	db *driver.EtcdDriver
	// application string
}

func NewInstaceRepository(driver *driver.EtcdDriver) *InstanceRepository {
	i := new(InstanceRepository)
	i.db = driver
	// i.application = app
	return i
}

func (i *InstanceRepository) Delete(id string) error {
	return nil
}

func (i *InstanceRepository) Get(id string) *model.InstanceModel {
	return nil
}

func (i *InstanceRepository) GetAll() []*model.InstanceModel {
	return nil
}

func (i *InstanceRepository) Set(entity *model.InstanceModel, app string) error {
	ops, err := ExtractEtcdOperations(entity)
	if err != nil {
		// TODO: Log Warning
		return err
	}
	for k, v := range ops {
		var t time.Time
		err := i.db.Set(k, false, v, t)
		if err != nil {
			// TODO: Log Warning
			return err
		}
	}
	return nil
}
