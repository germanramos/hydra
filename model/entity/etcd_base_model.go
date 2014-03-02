package entity

import (
	"errors"
	"reflect"
	"strconv"
)

type EtcdModelizer interface {
	ExportEtcdOperations() map[string]string
}

type EtcdBaseModel map[string]interface{}

func (e EtcdBaseModel) CastInterfaceToString(v interface{}) (string, error) {
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
		// TODO
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
			valueString, err := e.CastInterfaceToString(in)
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
