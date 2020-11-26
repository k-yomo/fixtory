package fixtory

import (
	"fmt"
	"reflect"
)

// MapNotZeroFields maps a struct fields to another struct fields if the field's value is not zero
func MapNotZeroFields[T any](from T, to *T) {
	fromKind := reflect.Indirect(reflect.ValueOf(from)).Type().Kind()
	if fromKind != reflect.Struct {
		panic(fmt.Sprintf("from must be struct, but got %s", fromKind))
	}
	toKind := reflect.Indirect(reflect.ValueOf(to)).Type().Kind()
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
