package main

import (
	"reflect"
	"testing"
)

func canBeNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return true
	default:
		return false
	}
}

func TestNewConfigShouldInitAllFields(t *testing.T) {
	config := newConfig()
	v := reflect.ValueOf(config)
	n := v.NumField()

	for i := 0; i < n; i++ {
		field := v.Field(i)
		if canBeNil(field) && field.IsNil() {
			tp := reflect.TypeOf(config)
			name := tp.Field(i).Name
			t.Errorf("Expected %s to not be nil", name)
		}
	}
}
