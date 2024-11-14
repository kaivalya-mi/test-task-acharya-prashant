package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func GetStructName(entity interface{}) string {
	if t := reflect.TypeOf(entity); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GetFieldFromStruct(entity interface{}, path string) (*reflect.Value, error) {
	scopes := strings.Split(path, ".")

	name := scopes[0]

	r := reflect.ValueOf(entity)
	f := reflect.Indirect(r).FieldByName(name)

	if len(scopes) > 1 {
		t := reflect.TypeOf(f)
		if t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("Not a struct %s", name)
		}
		return GetFieldFromStruct(f.Interface(), strings.Join(scopes[1:], "."))
	} else {
		return &f, nil
	}
}

func GetBoolFromStruct(entity interface{}, path string) bool {
	v, err := GetFieldFromStruct(entity, path)
	if err != nil {
		return false
	}
	return v.Bool()
}

func GetStringFromStruct(entity interface{}, path string) string {
	v, err := GetFieldFromStruct(entity, path)
	if err != nil {
		return ""
	}
	return v.String()
}

func GetIntFromStruct(entity interface{}, path string) int {
	return int(GetInt64FromStruct(entity, path))
}

func GetInt64FromStruct(entity interface{}, path string) int64 {
	v, err := GetFieldFromStruct(entity, path)
	if err != nil {
		return 0
	}
	return v.Int()
}

func GetUintFromStruct(entity interface{}, path string) uint {
	return uint(GetUint64FromStruct(entity, path))
}

func GetUint64FromStruct(entity interface{}, path string) uint64 {
	v, err := GetFieldFromStruct(entity, path)
	if err != nil {
		return 0
	}
	return v.Uint()
}

func GetFloatFromStruct(entity interface{}, path string) float64 {
	v, err := GetFieldFromStruct(entity, path)
	if err != nil {
		return 0
	}
	return v.Float()
}
