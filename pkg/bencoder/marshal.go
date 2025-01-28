package bencoder

import (
	"errors"
	"reflect"
)

func mapStructToMap(target interface{}) (map[string]interface{}, error) {
	reflectValue := reflect.ValueOf(target)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() != reflect.Struct {
		return nil, errors.New("struct or pointer to struct required")
	}

	mappedData := make(map[string]interface{})
	reflectType := reflectValue.Type()

	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		fieldValue := reflectValue.Field(i)

		bencodeTag, exists := field.Tag.Lookup("bencode")
		if !exists || !fieldValue.CanInterface() {
			continue
		}

		value, err := mapFieldValue(fieldValue)
		if err != nil {
			return nil, err
		}
		mappedData[bencodeTag] = value
	}

	return mappedData, nil
}

func mapFieldValue(fieldValue reflect.Value) (interface{}, error) {
	switch fieldValue.Kind() {
	case reflect.Struct:
		return mapStructToMap(fieldValue.Interface())
	case reflect.Ptr:
		if fieldValue.Elem().Kind() == reflect.Struct {
			return mapStructToMap(fieldValue.Interface())
		}
		return fieldValue.Interface(), nil
	case reflect.Slice, reflect.Array:
		return mapSliceOrArray(fieldValue)
	default:
		return fieldValue.Interface(), nil
	}
}

func mapSliceOrArray(sliceValue reflect.Value) ([]interface{}, error) {
	length := sliceValue.Len()
	result := make([]interface{}, 0, length)

	for i := 0; i < length; i++ {
		elemValue := sliceValue.Index(i)

		value, err := mapFieldValue(elemValue)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}
