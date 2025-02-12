package bencoder

import (
	"errors"
	"reflect"
)

type Bencoder interface {
	Decode(data []byte) (interface{}, error)
	Encode(data interface{}) ([]byte, error)
	Unmarshal(data []byte, target interface{}) error
	Marshal(target interface{}) ([]byte, error)
}

type SimpleBencoder struct{}

func NewSimpleBencoder() *SimpleBencoder {
	return &SimpleBencoder{}
}

func (bencoder *SimpleBencoder) Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	decodeFunc := getDecoder(data)
	return decodeFunc(data)
}

func (bencoder *SimpleBencoder) Encode(data interface{}) ([]byte, error) {
	return encode(data)
}

func (bencoder *SimpleBencoder) Unmarshal(data []byte, target interface{}) error {
	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return errors.New("struct type required for writing")
	}

	decodedData, err := bencoder.Decode(data)
	if err != nil {
		return err
	}

	mappedData, ok := decodedData.(map[string]interface{})
	if !ok {
		return errors.New("failed to cast decoded data to map")
	}

	for key, value := range mappedData {
		if err := setField(target, key, value); err != nil {
			return err
		}
	}

	return nil
}

func (bencoder *SimpleBencoder) Marshal(target interface{}) ([]byte, error) {
	mappedData, err := mapStructToMap(target)
	if err != nil {
		return nil, err
	}

	return bencoder.Encode(mappedData)
}
