package entity

import (
	"errors"
)

type Application struct {
	Id        string
	Balancers []Balancer
	Intances  []Instance
}

func NewApplication(id string, data map[string]interface{}) (Application, error) {
	balancers, err := extractBalancers(data)
	if err != nil {
		return nil, err
	}
	instances, err := extractInstances(data)
	if err != nil {
		return nil, err
	}
	return Application{
		Id:        id,
		Balancers: balancers,
		Instance:  instances,
	}
}

func NewApplicationFromEtcdBaseModel(e *EtcdBaseModel) (Application, error) {
	id, data, err := e.Explode()
	if err != nil {
		return nil, err
	}
	return NewApplication(id, data)
}

func checkIfDataContainsElementsInMap(data map[string]interface{}, key string) bool {
	mp, ok := data[key]
	if ok {
		mp, ok = mp.(map[string]interface{})
		if ok {
			if len(mp) > 0 {
				return true, nil
			} else {
				return false, errors.New("Incorrect type in " + key + ": expected map[string]interface{}")
			}
		}
	}
	return false, nil
}

func extractBalancers(data map[string]interface{}) ([]Balancer, error) {
	const BALANCERS_KEY string = "balancers"
	var balancers []Balancer = make([]Balancer, 0)
	hasBalancers, err := checkIfDataContainsElementsInMapKey(data, BALANCERS_KEY)
	if err != nil {
		return nil, err
	}
	if hasBalancers {
		return generateBalancers(data[BALANCERS_KEY])
	}
	return balancers, nil
}

func generateBalancers(balancers map[string]interface{}) ([]Balancers, error) {
	var balancerEntities []Balancer = make([]Balancer, 0)
	for id, data := range balancers {
		balancer, err := NewBalancer(id, data.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		balancerEntities = append(balancerEntities, balancer)
	}
	return balancerEntities, nil
}

func extractInstances(data map[string]interface{}) ([]Instance, error) {
	const INSTANCES_KEY string = "instances"
	var instances []Instance = make([]Instance, 0)
	hasInstances, err := checkIfDataContainsElementsInMapKey(data, INSTANCES_KEY)
	if err != nil {
		return nil, err
	}
	if hasInstances {
		return generateInstances(data[INSTANCES_KEY])
	}
	return instances, nil
}

func generateInstances(instances map[string]interface{}) ([]Instances, error) {
	var instanceEntities []Instance = make([]Instance, 0)
	for id, data := range instances {
		instance, err := NewInstance(id, data.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		instanceEntities = append(instanceEntities, instance)
	}
	return instanceEntities, nil
}
