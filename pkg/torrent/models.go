package torrent

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"reflect"
	"torrent/pkg/bencoder"
)

type TorrentFile struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	CreationDate int64      `bencode:"creation date"`
	Comment      string     `bencode:"comment"`
	CreatedBy    string     `bencode:"created by"`
	Info         InfoDict   `bencode:"info"`
}

type InfoDict struct {
	PieceLength int64  `bencode:"piece length"`
	Pieces      []byte `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int64  `bencode:"length"`
	Files       []File `bencode:"files"`
}

type File struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"`
}

func GeneratePieces(data []byte, pieceLength int) string {
	var buffer bytes.Buffer
	for i := 0; i < len(data); i += pieceLength {
		end := i + pieceLength
		if end > len(data) {
			end = len(data)
		}
		hash := sha1.Sum(data[i:end]) // TODO: make sure to use correct algorithm
		buffer.Write(hash[:])
	}
	return string(buffer.Bytes())
}

// TODO: move to bencoder package as marshal function
// TODO: refactor (case, smaller functions)

func NewTorrentFromBencode(data []byte) (*TorrentFile, error) {
	decodedData, err := bencoder.NewSimpleBencoder().Decode(data)
	if err != nil {
		return nil, err
	}

	mappedTorrent, ok := decodedData.(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to cast decoded data to map")
	}

	var torrentFile TorrentFile

	for key, value := range mappedTorrent {
		if err := setField(&torrentFile, key, value); err != nil {
			return nil, err
		}
	}

	return &torrentFile, nil
}

// setField sets the value of a struct field based on its bencode tag.
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("bencode")
		if tag != name {
			continue
		}
		fieldValue := structValue.Field(i)
		if !fieldValue.CanSet() {
			return errors.New("cannot set field " + name)
		}

		val := reflect.ValueOf(value)
		if fieldValue.Kind() == val.Kind() && fieldValue.Kind() != reflect.Slice { // slices needs special handling
			fieldValue.Set(val)
		} else if fieldValue.Kind() == reflect.String && val.Kind() == reflect.SliceOf(reflect.TypeOf(byte(1))).Kind() {
			fieldValue.Set(reflect.ValueOf(string(val.Bytes())))
		} else if fieldValue.Kind() == reflect.Struct && val.Kind() == reflect.Map {
			// Handle nested structs like InfoDict
			newStruct := reflect.New(fieldValue.Type()).Interface()
			for k, v := range value.(map[string]interface{}) {
				if err := setField(newStruct, k, v); err != nil {
					return err
				}
			}
			fieldValue.Set(reflect.ValueOf(newStruct).Elem())
		} else if fieldValue.Kind() == reflect.Slice && val.Kind() == reflect.Slice {
			// Use the new setSlice method to handle complex slice assignments
			if err := setSlice(fieldValue, val); err != nil {
				return err
			}
		} else {
			return errors.New("type mismatch for field " + name)
		}
		return nil
	}
	return nil
}

func setSlice(target reflect.Value, value reflect.Value) error {
	if target.Kind() != reflect.Slice || value.Kind() != reflect.Slice {
		return fmt.Errorf("target and value must be slices")
	}

	// Create a new slice with the same type and length as the target
	newSlice := reflect.MakeSlice(target.Type(), value.Len(), value.Len())

	for i := 0; i < value.Len(); i++ {
		elem := reflect.ValueOf(value.Index(i).Interface())
		targetElem := newSlice.Index(i)

		// Check if the element is a slice itself
		if (target.Type().Elem().Kind() == reflect.Slice) && (elem.Kind() == reflect.Slice) {
			// Recursively handle slices of slices
			if err := setSlice(targetElem, elem); err != nil {
				return err
			}
		} else if target.Type().Elem().Kind() == reflect.String && elem.Kind() == reflect.SliceOf(reflect.TypeOf(byte(1))).Kind() {
			targetElem.Set(reflect.ValueOf(string(elem.Bytes())))
		} else if target.Type().Elem().Kind() == reflect.Int64 && elem.Kind() == reflect.Int64 {
			targetElem.Set(reflect.ValueOf(elem))
		} else if target.Type().Elem().Kind() == reflect.Struct && elem.Kind() == reflect.Map {
			// Recursively handle slices of maps converting them to structs
			newElem := reflect.New(target.Type().Elem()).Elem()
			for k, v := range elem.Interface().(map[string]interface{}) {
				if err := setField(newElem.Addr().Interface(), k, v); err != nil {
					return err
				}
			}
			targetElem.Set(newElem)
		} else {
			// Handle other cases directly
			targetElem.Set(elem)
		}
	}

	// Set the final slice value to the target
	target.Set(newSlice)
	return nil
}
