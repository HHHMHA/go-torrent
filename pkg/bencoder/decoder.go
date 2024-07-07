package bencoder

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"reflect"
	"strconv"
)

func getDecoder(data []byte) func([]byte) (interface{}, error) {
	var decodeFuncs = map[byte]func([]byte) (interface{}, error){
		'i': decodeInt,
		'l': decodeList,
	}

	if decodeFunc, ok := decodeFuncs[data[0]]; ok {
		return decodeFunc
	}
	return decodeString
}

func decodeIntInList(data []byte, startIndex int) (int, interface{}, error) {
	if len(data) <= 3 || data[startIndex] != 'i' {
		return 0, nil, fmt.Errorf("invalid integer format: %s", string(data))
	}
	// find the endIndex
	endIndex := startIndex + 1
	for data[endIndex] != byte('e') {
		endIndex += 1
		if endIndex > len(data) {
			// WTF why are you passing this asshole it's not an int
			return 0, nil, fmt.Errorf("invalid integer format: %s", string(data))
		}
	}
	result, err := decodeInt(data[startIndex : endIndex+1])
	return endIndex + 1, result, err
}

func decodeInt(data []byte) (interface{}, error) {
	if len(data) < 3 || data[0] != 'i' || data[len(data)-1] != 'e' {
		return nil, fmt.Errorf("invalid integer format: %s", string(data))
	}
	return strconv.ParseInt(string(data[1:len(data)-1]), 10, 64)
}

func decodeStringInList(data []byte) (int, []byte, error) {
	// It's byte string not normal string so we will return the bytes
	// Very handy too since the piece hash is in bytes
	separator := []byte(":")
	result := bytes.Split(data, separator)
	if len(result) < 2 {
		return 0, nil, errors.New("unknown format")
	}

	length, err := strconv.Atoi(string(result[0]))
	if err != nil {
		return 0, nil, errors.New("length of string is not correct")
	}

	// TODO: Assuming ascii for now, handle UTF-8 later
	if length > len(result[1]) {
		return 0, nil, errors.New("mismatch of length and byte string")
	}

	return length, result[1], nil
}

func decodeString(data []byte) (interface{}, error) {
	length, result, err := decodeStringInList(data)
	if err != nil {
		return nil, err
	}

	// TODO: Assuming ascii for now, handle UTF-8 later
	if length != len(result) {
		return nil, errors.New("mismatch of length and byte string")
	}

	return result, nil
}

func decodeListInner(data []byte, startIndex int) (nextIndex int, result interface{}, err error) {
	result = []interface{}{}
	err = fmt.Errorf("list element format invalid")
	nextIndex = startIndex + 1 // skip the first l

	if len(data) < nextIndex || data[startIndex] != 'l' {
		return
	}

	var element any
	var innerErr error

	for data[nextIndex] != byte('e') {
		decoderFunc := getDecoder(data[nextIndex : nextIndex+1])
		switch reflect.ValueOf(decoderFunc).Pointer() {
		case reflect.ValueOf(decodeInt).Pointer():
			nextIndex, element, innerErr = decodeIntInList(data, nextIndex)
			if innerErr != nil {
				result = nil

				return
			}
		case reflect.ValueOf(decodeString).Pointer():
			currentSlice := data[nextIndex:]
			var length int
			length, element, innerErr = decodeStringInList(currentSlice)
			if innerErr != nil {
				result = nil
				return
			}

			colonIndex := slices.IndexFunc(currentSlice, func(e byte) bool { return e == byte(':') })
			element = element.([]uint8)[:length]
			nextIndex += colonIndex + length + 1
		case reflect.ValueOf(decodeList).Pointer():
			nextIndex, element, innerErr = decodeListInner(data, nextIndex)
			if nextIndex < len(data)-1 {
				nextIndex += 1 // skip the e for next decode
			}
		}
		result = append(result.([]interface{}), element)
		if nextIndex >= len(data) {
			result = nil
			err = fmt.Errorf("list element format invalid")
			return
		}
	}

	return
}

func decodeList(data []byte) (interface{}, error) {
	_, result, err := decodeListInner(data, 0)
	return result, err
}
