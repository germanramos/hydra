package entity

type Instance struct {
	Id   string
	Info map[string]interface{}
}

func NewInstance(id string, data map[string]interface{}) (Application, error) {
	return Balancer{
		Id:   id,
		Info: data,
	}, nil
}

func NewBalancerFromEtcdBaseModel(e *EtcdBaseModel) (Application, error) {
	id, data, err := e.Explode()
	if err != nil {
		return nil, err
	}
	return NewBalancer(id, data)
}
