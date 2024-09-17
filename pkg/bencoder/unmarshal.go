package bencoder

import (
	"errors"
	"fmt"
	"reflect"
)

// setField sets the value of a struct field based on its bencode tag.
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("bencode")
		if tag != name {
			continue
		}
		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			return errors.New("cannot set field " + name)
		}

		val := reflect.ValueOf(value)
		if err := assignValue(fieldValue, val); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// assignValue handles the assignment of values based on their types.
func assignValue(target reflect.Value, value reflect.Value) error {
	if target.Kind() == value.Kind() && target.Kind() != reflect.Slice {
		target.Set(value)
		return nil
	}

	// Handle specific type conversions
	switch target.Kind() {
	case reflect.String:
		if isByteSlice(value) {
			target.SetString(string(value.Bytes()))
			return nil
		}
	case reflect.Struct:
		if value.Kind() == reflect.Map {
			return assignStruct(target, value)
		}
	case reflect.Slice:
		if value.Kind() == reflect.Slice {
			return assignSlice(target, value)
		}
	}

	return fmt.Errorf("type mismatch for field %v", target.Type().Name())
}

// isByteSlice checks if a reflect.Value is a slice of bytes.
func isByteSlice(value reflect.Value) bool {
	return value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Uint8
}

// assignStruct assigns a map's values to a struct's fields recursively.
func assignStruct(target reflect.Value, value reflect.Value) error {
	newStruct := reflect.New(target.Type()).Interface()
	for k, v := range value.Interface().(map[string]interface{}) {
		if err := setField(newStruct, k, v); err != nil {
			return err
		}
	}
	target.Set(reflect.ValueOf(newStruct).Elem())
	return nil
}

// assignSlice handles complex slice assignments.
func assignSlice(target reflect.Value, value reflect.Value) error {
	newSlice := reflect.MakeSlice(target.Type(), value.Len(), value.Len())

	for i := 0; i < value.Len(); i++ {
		elem := reflect.ValueOf(value.Index(i).Interface())
		targetElem := newSlice.Index(i)

		if targetElem.Kind() == reflect.Slice && elem.Kind() == reflect.Slice {
			// Recursively handle slices of slices
			if err := assignSlice(targetElem, elem); err != nil {
				return err
			}
		} else if targetElem.Kind() == reflect.Struct && elem.Kind() == reflect.Map {
			// Convert map to struct (for cases like []File)
			if err := assignStruct(targetElem, elem); err != nil {
				return err
			}
		} else {
			if err := assignValue(targetElem, elem); err != nil {
				return err
			}
		}
	}

	target.Set(newSlice)
	return nil
}
