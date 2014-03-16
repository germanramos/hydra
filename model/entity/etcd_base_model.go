package entity

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/store"
)

type EtcdModelizer interface {
	ExportEtcdOperations() map[string]string
}

type EtcdBaseModel map[string]interface{}

type EtcdBaseModels []EtcdBaseModel

func NewModelFromEvent(event *store.Event) (*EtcdBaseModel, error) {
	model := make(map[string]interface{})
	if err := proccessStruct(event, model); err != nil {
		return nil, err
	}
	m := EtcdBaseModel(model)
	return &m, nil
}

func NewModelsFromEvent(event *store.Event) (*EtcdBaseModels, error) {
	models := make([]EtcdBaseModel, 0)
	nodes := []*store.NodeExtern(reflect.ValueOf(event).Elem().FieldByName("Node").Elem().FieldByName("Nodes").Interface().(store.NodeExterns))
	for _, node := range nodes {
		model := make(map[string]interface{})
		if err := proccessStruct(node, model); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	m := EtcdBaseModels(models)
	return &m, nil
}

func ExtractJsonKeyFromEtcdKey(s string) (string, error) {
	lastIndex := strings.LastIndex(s, "/")
	if lastIndex == -1 {
		// TODO
		return "", errors.New("Bad etcd key")
	}
	key := s[lastIndex+1:]
	if len(key) == 0 {
		// TODO
		return "", errors.New("Bad etcd key")
	}
	return key, nil
}

func proccessStruct(s interface{}, m map[string]interface{}) error {
	if node := reflect.ValueOf(s).Elem().FieldByName("Node"); node.IsValid() {
		proccessStruct(node.Interface(), m)
	} else if exists := reflect.ValueOf(s).Elem().FieldByName("Nodes"); exists.IsValid() && !exists.IsNil() {
		key, _ := ExtractJsonKeyFromEtcdKey(reflect.ValueOf(s).Elem().FieldByName("Key").Interface().(string))
		nodes := []*store.NodeExtern(reflect.ValueOf(s).Elem().FieldByName("Nodes").Interface().(store.NodeExterns))
		m[key] = make(map[string]interface{})
		for _, node := range nodes {
			proccessStruct(node, m[key].(map[string]interface{}))
		}
	} else {
		value, err := CastInterfaceToString(reflect.ValueOf(s).Elem().FieldByName("Value").Interface())
		if err != nil {
			return err
		}
		key, _ := ExtractJsonKeyFromEtcdKey(reflect.ValueOf(s).Elem().FieldByName("Key").Interface().(string))
		m[key] = value
	}
	return nil
}

// func (e EtcdBaseModel) CastInterfaceToString(v interface{}) (string, error) {
func CastInterfaceToString(v interface{}) (string, error) {
	var str string
	switch v.(type) {
	case nil:
		str = ""
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = strconv.Itoa(v.(int))
	case float32, float64:
		str = strconv.FormatFloat(v.(float64), 'f', 2, 64)
	case bool:
		str = strconv.FormatBool(v.(bool))
	case string:
		str = v.(string)
	default:
		// TODO: improve error
		return "", errors.New("Bad interface")
	}
	return str, nil
}

func (e EtcdBaseModel) ExportEtcdOperations() (map[string]string, error) {
	var operations map[string]string
	operations = make(map[string]string)

	var processInterface func(interface{}, string) error
	var processMap func(map[string]interface{}, string) error
	var processSlice func([]interface{}, string) error

	processInterface = func(in interface{}, key string) error {
		switch reflect.ValueOf(in).Kind() {
		case reflect.Map:
			processMap(in.(map[string]interface{}), key)
		case reflect.Slice:
			processSlice(in.([]interface{}), key)
		default:
			valueString, err := CastInterfaceToString(in)
			if err != nil {
				return err
			}
			operations[key] = valueString
		}
		return nil
	}

	processSlice = func(s []interface{}, parentKey string) error {
		for key, value := range s {
			if err := processInterface(value, parentKey+"/"+strconv.Itoa(key)); err != nil {
				return err
			}
		}
		return nil
	}

	processMap = func(mp map[string]interface{}, parentKey string) error {
		for key, value := range mp {
			if err := processInterface(value, parentKey+"/"+key); err != nil {
				return err
			}
		}
		return nil
	}

	if err := processMap(e, ""); err != nil {
		return nil, err
	}
	return operations, nil
}

func (e *EtcdBaseModel) Explode() (string, map[string]interface{}) {
	eP := map[string]interface{}(*e)
	if len(eP) != 1 {
		return "", nil
	}
	for key, value := range eP {
		return key, value.(map[string]interface{})
	}
	return "", nil
}
