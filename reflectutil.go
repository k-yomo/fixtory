package fixtory

import (
	"fmt"
	"reflect"
)

// MapNotZeroFields maps a struct fields to another struct fields if the field's value is not zero
// from and to must be same struct, and to must be pointer
func MapNotZeroFields(from interface{}, to interface{}) {
	fromKind := reflect.Indirect(reflect.ValueOf(from)).Type().Kind()
	if fromKind != reflect.Struct {
		panic(fmt.Sprintf("from must be struct, but got %s", fromKind))
	}
	toKind :=reflect.Indirect(reflect.ValueOf(to)).Type().Kind()
	if toKind != reflect.Struct {
		panic(fmt.Sprintf("to must be struct, but got %s", toKind))
	}

	fromV := reflect.Indirect(reflect.ValueOf(from))
	toV := reflect.ValueOf(to).Elem()

	for i := 0; i < fromV.NumField(); i++ {
		fieldV := fromV.Field(i)
		if !fieldV.IsZero() {
			toV.FieldByName(fromV.Type().Field(i).Name).Set(fieldV)
		}
	}
}

// ConvertToInterfaceArray converts any type of array to interface array
func ConvertToInterfaceArray(from interface{}) []interface{} {
	fromType := reflect.TypeOf(from)
	kind := fromType.Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		panic(fmt.Sprintf("from must be array or slice, but got %s", kind))
	}

	fromValue := reflect.ValueOf(from)
	interfaceArray := make([]interface{}, 0, fromValue.Len())
	for i := 0; i < fromValue.Len(); i++ {
		interfaceArray = append(interfaceArray, fromValue.Index(i).Interface())
	}
	return interfaceArray
}
