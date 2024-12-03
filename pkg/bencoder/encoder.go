package bencoder

import (
	"errors"
	"fmt"
	"reflect"
)

func encode(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, errors.New("no data to encode")
	}
	dataType := reflect.TypeOf(data).Kind()
	switch dataType {
	case reflect.Int, reflect.Int64:
		return encodeInt(data)
	case reflect.String:
		return encodeString(data)
	case reflect.Slice:
		if reflect.TypeOf(data) == reflect.TypeOf([]byte(nil)) {
			return encodeString(data)
		}
		return encodeList(data)
	case reflect.Array:
		return encodeList(data)
	}

	return nil, errors.New("unsupported type")
}

func encodeInt(data interface{}) ([]byte, error) {
	value, ok := data.(int64)
	if !ok {
		return nil, errors.New("invalid int provided")
	}
	return []byte(fmt.Sprintf("i%de", value)), nil
}

func encodeString(data interface{}) ([]byte, error) {
	var valueArray []uint8
	var value string
	ok := true
	if reflect.TypeOf(data).Kind() != reflect.String {
		valueArray, ok = data.([]uint8)
		if !ok {
			return nil, errors.New("invalid string provided")
		}
		value = string(valueArray)
	} else {
		value, ok = data.(string)
	}

	if !ok {
		return nil, errors.New("invalid string provided")
	}
	return []byte(fmt.Sprintf("%d:%s", len(value), value)), nil
}

func encodeList(data interface{}) ([]byte, error) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return nil, errors.New("input data is not a list")
	}
	var result []byte
	result = append(result, 'l') // Start with 'l'
	s := reflect.ValueOf(data)

	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i).Interface()
		encodedElem, err := encode(elem)
		if err != nil {
			return nil, err
		}
		result = append(result, encodedElem...)
	}
	result = append(result, 'e') // End with 'e'
	return result, nil
}
