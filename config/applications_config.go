package config

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/innotech/hydra/model/entity"
	. "github.com/innotech/hydra/model/repository"
)

type ApplicationsConfig struct {
	Apps EtcdBaseModels
	Repo *EtcdBaseRepository
}

func NewApplicationsConfig() *ApplicationsConfig {
	a := new(ApplicationsConfig)
	a.Repo = NewEctdRepository()
	a.Repo.SetCollection("applications")
	return a
}

func (a *ApplicationsConfig) Load(pathToConfigFile string) error {
	if err := a.loadAppsFromJSON(pathToConfigFile); err != nil {
		return err
	}
	return nil
}

func (a *ApplicationsConfig) loadAppsFromJSON(pathToFile string) error {
	fileContent, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(fileContent, &(a.Apps)); err != nil {
		return err
	}
	return nil
}

// TODO: Test
func (a *ApplicationsConfig) Persists() error {
	for _, app := range a.Apps {
		err := a.Repo.Set(&app)
		if err != nil {
			// TODO: delete applications directory
			return err
		}
	}
	return nil
}
