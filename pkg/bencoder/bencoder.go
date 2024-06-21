package bencoder

import "errors"

type Bencoder interface {
	Decode(data []byte) (interface{}, error)
	Encode(data interface{}) ([]byte, error)
}

type SimpleBencoder struct{}

func NewSimpleBencoder() *SimpleBencoder {
	return &SimpleBencoder{}
}

func (bencoder *SimpleBencoder) Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	if decodeFunc, ok := decodeFuncs[data[0]]; ok {
		return decodeFunc(data)
	}

	return nil, errors.New("invalid data")
}

func (bencoder *SimpleBencoder) Encode(data interface{}) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
