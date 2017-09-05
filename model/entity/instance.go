package entity

type Instance struct {
	Id   string
	Info map[string]interface{}
}

func NewInstance(id string, data map[string]interface{}) (Instance, error) {
	return Instance{
		Id:   id,
		Info: data,
	}, nil
}

func NewInstanceFromEtcdBaseModel(e *EtcdBaseModel) (Instance, error) {
	id, data, err := e.Explode()
	if err != nil {
		return Instance{}, err
	}
	return NewInstance(id, data)
}
