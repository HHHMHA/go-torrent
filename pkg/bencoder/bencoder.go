package bencoder

import (
	"errors"
	"reflect"
)

type Bencoder interface {
	Decode(data []byte) (interface{}, error)
	Encode(data interface{}) ([]byte, error)
	Unmarshal(data []byte, target interface{}) error
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
	//TODO implement me
	panic("implement me")
}

func (bencoder *SimpleBencoder) Unmarshal(data []byte, target interface{}) error {
	if reflect.TypeOf(target).Elem().Kind() != reflect.Struct {
		return errors.New("struct type required for writing")
	}

	decodedData, err := bencoder.Decode(data)
	if err != nil {
		return err
	}

	mappedTorrent, ok := decodedData.(map[string]interface{})
	if !ok {
		return errors.New("failed to cast decoded data to map")
	}

	for key, value := range mappedTorrent {
		if err := setField(target, key, value); err != nil {
			return err
		}
	}

	return nil
}
