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
