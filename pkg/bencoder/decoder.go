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
		'd': decodeDict,
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

func decodeString(data []byte) (interface{}, error) {
	_, result, err := decodeStringInList(data, 0)
	if err != nil {
		return nil, err
	}

	// TODO: Assuming ascii for now, handle UTF-8 later
	lengthPart, err := strconv.Atoi(string(data[0 : len(data)-len(result)-1]))
	if err != nil || lengthPart != len(result) {
		return nil, errors.New("mismatch of length and byte string")
	}

	return result, nil
}

func decodeStringInList(data []byte, startIndex int) (int, []byte, error) {
	// It's byte string not normal string so we will return the bytes
	// Very handy too since the piece hash is in bytes
	separator := []byte(":")
	lengthBytes, stringBytes, found := bytes.Cut(data, separator)
	if !found {
		return 0, nil, errors.New("unknown format")
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return 0, nil, errors.New("length of string is not correct")
	}

	// TODO: Assuming ascii for now, handle UTF-8 later
	if length > len(stringBytes) {
		return 0, nil, errors.New("mismatch of length and byte string")
	}

	colonIndex := slices.IndexFunc(data, func(e byte) bool { return e == byte(':') })
	stringBytes = stringBytes[:length]
	nextIndex := startIndex + colonIndex + length + 1
	return nextIndex, stringBytes, nil
}

func decodeList(data []byte) (interface{}, error) {
	_, result, err := decodeListInner(data, 0)
	return result, err
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
			nextIndex, element, innerErr = decodeStringInList(data[nextIndex:], nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
		case reflect.ValueOf(decodeList).Pointer():
			nextIndex, element, innerErr = decodeListInner(data, nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
			if nextIndex < len(data)-1 {
				nextIndex += 1 // skip the e for next decode
			}
		case reflect.ValueOf(decodeDict).Pointer():
			nextIndex, element, innerErr = decodeDictInner(data, nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
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

	err = nil
	return
}

func decodeDict(data []byte) (interface{}, error) {
	_, result, err := decodeDictInner(data, 0)
	return result, err
}

func decodeDictInner(data []byte, startIndex int) (nextIndex int, result interface{}, err error) {
	result = map[string]interface{}{}
	err = fmt.Errorf("invalid dictionary format")
	nextIndex = startIndex + 1 // skip the first d

	if len(data) < nextIndex || data[startIndex] != 'd' {
		return
	}

	var element any
	var key []byte
	var innerErr error

	for data[nextIndex] != byte('e') {
		// first extract the key and move the index after the extracted key string
		nextIndex, key, innerErr = decodeStringInList(data[nextIndex:], nextIndex)
		if innerErr != nil || len(data) < nextIndex {
			result = nil
			return
		}
		if _, exists := result.(map[string]interface{})[string(key)]; exists {
			result = nil
			return
		}
		// now get the value
		decoderFunc := getDecoder(data[nextIndex : nextIndex+1])
		switch reflect.ValueOf(decoderFunc).Pointer() {
		case reflect.ValueOf(decodeInt).Pointer():
			nextIndex, element, innerErr = decodeIntInList(data, nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
		case reflect.ValueOf(decodeString).Pointer():
			nextIndex, element, innerErr = decodeStringInList(data[nextIndex:], nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
		case reflect.ValueOf(decodeList).Pointer():
			nextIndex, element, innerErr = decodeListInner(data, nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
			if nextIndex < len(data)-1 {
				nextIndex += 1 // skip the e for next decode
			}
		case reflect.ValueOf(decodeDict).Pointer():
			nextIndex, element, innerErr = decodeDictInner(data, nextIndex)
			if innerErr != nil {
				result = nil
				return
			}
			if nextIndex < len(data)-1 {
				nextIndex += 1 // skip the e for next decode
			}
		}
		result.(map[string]interface{})[string(key)] = element
		if nextIndex >= len(data) {
			result = nil
			return
		}
	}

	err = nil
	return
}
