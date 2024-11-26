package bencoder

import (
	"errors"
	"fmt"
	"reflect"
)

func decode(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, errors.New("no data to encode")
	}
	dataType := reflect.TypeOf(data).Kind()
	switch dataType {
	case reflect.Int, reflect.Int64:
		return encodeInt(data)
	case reflect.String, reflect.SliceOf(reflect.TypeOf(byte('a'))).Kind(): // TODO: probably should change and make it normal list?
		return encodeString(data)
	}

	return nil, nil
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
