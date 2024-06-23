package bencoder

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var decodeFuncs = map[byte]func([]byte) (interface{}, error){
	'i': decodeInt,
}

func decodeInt(data []byte) (interface{}, error) {
	if len(data) < 3 || data[0] != 'i' || data[len(data)-1] != 'e' {
		return nil, fmt.Errorf("invalid integer format: %s", string(data))
	}
	return strconv.ParseInt(string(data[1:len(data)-1]), 10, 64)
}

func decodeString(data []byte) (interface{}, error) {
	// It's byte string not normal string so we will return the bytes
	// Very handy too since the piece hash is in bytes
	separator := []byte(":")
	result := bytes.Split(data, separator)
	if len(result) != 2 {
		return nil, errors.New("unknown format")
	}

	length, err := strconv.Atoi(string(result[0]))
	if err != nil {
		return nil, errors.New("length of string is not correct")
	}

	// TODO: Assuming ascii for now, handle UTF-8 later
	if length != len(result[1]) {
		return nil, errors.New("mismatch of length and byte string")
	}

	return result[1], nil
}
