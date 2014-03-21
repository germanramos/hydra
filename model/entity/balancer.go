package entity

type Balancer struct {
	Id   string
	Args map[string]interface{}
}

func NewBalancer(id string, data map[string]interface{}) (Balancer, error) {
	return Balancer{
		Id:   id,
		Args: data,
	}, nil
}

func NewBalancerFromEtcdBaseModel(e *EtcdBaseModel) (Balancer, error) {
	id, data, err := e.Explode()
	if err != nil {
		return Balancer{}, err
	}
	return NewBalancer(id, data)
}
