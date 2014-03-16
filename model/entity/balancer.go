package entity

type Balancer struct {
	Id   string
	Args map[string]interface{}
}

func NewBalancer(id string, data map[string]interface{}) (Application, error) {
	return Balancer{
		Id:   id,
		Args: data,
	}
}

func NewBalancerFromEtcdBaseModel(e *EtcdBaseModel) (Application, error) {
	id, data, err := e.Explode()
	if err != nil {
		return nil, err
	}
	return NewBalancer(id, data)
}
