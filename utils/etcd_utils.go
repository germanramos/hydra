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
		kind := reflect.TypeOf(interfc).Kind()
		if kind == reflect.Map {
			fmt.Println("********** MAP *********")
			// jsonObject := reflect.ValueOf(interfc).Elem().Interface().(map[string]interface{})
			// TRY:
			jsonObject := interfc.(map[string]interface{})
			for key, value := range jsonObject {
				kind := reflect.ValueOf(value).Kind()
				if kind == reflect.Struct || kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
					fmt.Println("Entra: " + parentKey + "/" + key)
					err := processInterface(value, parentKey+"/"+key)
					if err != nil {
						return err
					}
				} else {
					value, err := CastInterfaceToString(value)
					if err != nil {
						return err
					}
					operations[parentKey+"/"+key] = value
				}
				// switch value.(type) {
				// case interface{}:
				// 	processInterface(value, parentKey+"/"+key)
				// default:
				// 	value, err := CastInterfaceToString(value)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	operations[parentKey+"/"+key] = value
				// }
			}
		} else if kind == reflect.Struct {
			numOfFields := reflect.ValueOf(interfc).NumField()
			fieldIndex := []int{0}
			for i := 0; i < numOfFields; i++ {
				fieldIndex[0] = i
				var structField reflect.StructField
				structField = reflect.TypeOf(interfc).FieldByIndex(fieldIndex)
				key := parentKey + "/" + structField.Name
				var structValue reflect.Value
				structValue = reflect.ValueOf(interfc).FieldByIndex(fieldIndex)
				// switch structValue.Interface().(type) {
				// case interface{}:
				// 	fmt.Println("Entra: " + key)
				// 	err := processInterface(structValue.Interface(), key)
				// 	if err != nil {
				// 		return err
				// 	}
				// default:
				// 	value, err := CastInterfaceToString(structValue.Interface())
				// 	if err != nil {
				// 		return err
				// 	}
				// 	operations[key] = value
				// }
				fmt.Println("Entra: " + key)
				if structValue.Kind() == reflect.Struct || structValue.Kind() == reflect.Map || structValue.Kind() == reflect.Array || structValue.Kind() == reflect.Slice {
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
		} else if kind == reflect.Slice || kind == reflect.Array {

			// TODO: Test Array type
			// jsonArray := interfc.([]interface{})
			// jsonArray := reflect.ValueOf(interfc).Elem().Slice(i, j)
			jsonArrayValue := reflect.ValueOf(interfc)
			length := jsonArrayValue.Len()
			fmt.Println("Entra Len:" + strconv.Itoa(length))
			for i := 0; i < length; i++ {
				value := jsonArrayValue.Index(i).Interface()
				kind := jsonArrayValue.Index(i).Kind()
				fmt.Println(kind.String())
				if kind == reflect.Interface || kind == reflect.Struct || kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
					fmt.Println("Entra: " + parentKey + "/" + strconv.Itoa(i))
					err := processInterface(value, parentKey+"/"+strconv.Itoa(i))
					if err != nil {
						return err
					}
				} else {
					fmt.Println("Entra 2: " + parentKey + "/" + strconv.Itoa(i))
					value, err := CastInterfaceToString(value)
					if err != nil {
						return err
					}
					operations[parentKey+"/"+strconv.Itoa(i)] = value
				}
				// switch value.(type) {
				// case interface{}:
				// 	processInterface(value, parentKey+"/"+strconv.Itoa(i))
				// default:
				// 	value, err := CastInterfaceToString(value)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	operations[parentKey+"/"+strconv.Itoa(i)] = value
				// }
			}
		}
		// 	// jsonArray := interfc.([]interface{})
		// 	// for key, value := range jsonArray {
		// 	// 	fmt.Printf("param #%d is an int\n", key)
		// 	// 	switch value.(type) {
		// 	// 	case interface{}:
		// 	// 		processInterface(value, parentKey+"/"+strconv.Itoa(key))
		// 	// 	default:
		// 	// 		value, err := CastInterfaceToString(value)
		// 	// 		if err != nil {
		// 	// 			return err
		// 	// 		}
		// 	// 		operations[parentKey+"/"+strconv.Itoa(key)] = value
		// 	// 	}
		// 	// }
		// }
		return nil

		// numOfFields := 0
		// numOfFields = reflect.ValueOf(interfc).NumField()
		// fieldIndex := []int{0}
		// for i := 0; i < numOfFields; i++ {
		// 	fieldIndex[0] = i
		// 	var structField reflect.StructField
		// 	structField = reflect.TypeOf(interfc).FieldByIndex(fieldIndex)
		// 	key := parentKey + "/" + structField.Name
		// 	var structValue reflect.Value
		// 	structValue = reflect.ValueOf(interfc).FieldByIndex(fieldIndex)
		// 	if structValue.Kind() == reflect.Struct {
		// 		err := processInterface(structValue.Interface(), key)
		// 		if err != nil {
		// 			return err
		// 		}
		// 	} else {
		// 		value, err := CastInterfaceToString(structValue.Interface())
		// 		if err != nil {
		// 			return err
		// 		}
		// 		operations[key] = value
		// 	}
		// }
		// return nil
	}

	typ := reflect.TypeOf(interfc)
	if typ.Kind() == reflect.Ptr && (typ.Elem().Kind() == reflect.Struct || typ.Elem().Kind() == reflect.Map) {
		err := processInterface(reflect.ValueOf(interfc).Elem().Interface(), "")
		if err != nil {
			return nil, err
		}
		return operations, nil
	} else {
		// TODO:
		panic("Bad pointer type has been passed")
	}
}

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
		return "", errors.New("Bad interface")
	}
	return str, nil
}
