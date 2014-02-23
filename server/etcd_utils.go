package server

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func DecodeJsonObject(data []byte, v interface{}) (map[string]string, error) {
	err := json.Unmarshal(data, v)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("%+v", v)
	// var ops map[string]string
	ops, err := extractEtcdOperations(v)

	return ops, err
}

func extractEtcdOperations(v interface{}) (map[string]string, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

}
