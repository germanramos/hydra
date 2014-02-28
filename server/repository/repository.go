package repository

type Repository interface {
	Delete(id string) error
	Get(id string) interface{}
	GetAll() []interface{}
	Set(entity interface{}) error
}
