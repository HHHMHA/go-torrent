package bencoder

import (
	"fmt"
	"strconv"
)

/**
This package focuses on speed above all instead of careful error handling or readability
So you will notice we won't be usig Atoi for example to avoid overhead of creating a string object
We won't validate the input and instead just throw a normal error if something happens
*/

var decodeFuncs = map[byte]func([]byte) (interface{}, error){
	'i': decodeInt,
}

func decodeInt(data []byte) (interface{}, error) {
	if len(data) < 3 || data[0] != 'i' || data[len(data)-1] != 'e' {
		return nil, fmt.Errorf("invalid integer format: %s", string(data))
	}
	return strconv.ParseInt(string(data[1:len(data)-1]), 10, 64)
}
