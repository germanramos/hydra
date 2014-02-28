package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func DecodeJsonObject(data []byte, v interface{}) (map[string]string, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		// TODO: move out
		fmt.Println("error:", err)
	}
	// ops, err := extractEtcdOperations(v, false)
	ops, err := ExtractEtcdOperations(v)

	return ops, err
}

func ExtractEtcdOperations(interfc interface{} /*, whole bool*/) (map[string]string, error) {
	var operations map[string]string
	operations = make(map[string]string)

	// TODO: Change name to processInterface
	var processInterface func(interface{}, string) error
	processInterface = func(interfc interface{}, parentKey string) error {
		numOfFields := 0
		// if reflect.TypeOf(interfc).Kind() == reflect.Ptr {
		// 	numOfFields = reflect.TypeOf(interfc).Elem().NumField()
		// } else {
		numOfFields = reflect.ValueOf(interfc).NumField()
		// }
		fieldIndex := []int{0}
		for i := 0; i < numOfFields; i++ {
			fieldIndex[0] = i
			var structField reflect.StructField
			// if reflect.TypeOf(interfc).Kind() == reflect.Ptr {
			// 	structField = reflect.TypeOf(interfc).Elem().FieldByIndex(fieldIndex)
			// } else {
			structField = reflect.TypeOf(interfc).FieldByIndex(fieldIndex)
			// }
			key := parentKey + "/" + structField.Name
			var structValue reflect.Value
			// if reflect.TypeOf(interfc).Kind() == reflect.Ptr {
			// 	structValue = reflect.ValueOf(interfc).Elem().FieldByIndex(fieldIndex)
			// } else {
			structValue = reflect.ValueOf(interfc).FieldByIndex(fieldIndex)
			// }
			if structValue.Kind() == reflect.Struct {
				/*if whole == true {
					operations[key] = ""
				}*/
				err := processInterface(structValue.Interface(), key)
				if err != nil {
					return err
				}
			} else {
				value, err := CastInterfaceToString(structValue.Interface())
				if err != nil {
					return err
				}
				operations[key] = value
			}
		}
		return nil
	}

	// TODO: Arrays
	typ := reflect.TypeOf(interfc)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		err := processInterface(reflect.ValueOf(interfc).Elem().Interface(), "")
		if err != nil {
			return nil, err
		}
		return operations, nil
	} else {
		return nil, errors.New("Bad pointer")
	}
}

func CastInterfaceToString(v interface{}) (string, error) {
	var str string
	switch v.(type) {
	case int:
		str = strconv.Itoa(v.(int))
	case int8:
		str = strconv.Itoa(v.(int))
	case int16:
		str = strconv.Itoa(v.(int))
	case int32:
		str = strconv.Itoa(v.(int))
	case int64:
		str = strconv.Itoa(v.(int))
	case uint:
		str = strconv.Itoa(v.(int))
	case uint8:
		str = strconv.Itoa(v.(int))
	case uint16:
		str = strconv.Itoa(v.(int))
	case uint32:
		str = strconv.Itoa(v.(int))
	case uint64:
		str = strconv.Itoa(v.(int))
	case float32:
		str = strconv.FormatFloat(v.(float64), 'f', 2, 32)
	case float64:
		str = strconv.FormatFloat(v.(float64), 'f', 2, 64)
	case bool:
		str = strconv.FormatBool(v.(bool))
	case string:
		str = v.(string)
	default:
		return "", errors.New("Bad interface")
	}
	return str, nil
}
