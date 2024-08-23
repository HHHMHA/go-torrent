package torrent

import (
	"bytes"
	"crypto/sha1"
	"errors"
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
			// Handle slices of slices (AnnounceList) and slices of structs (Files)
			if fieldValue.Type().Elem().Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Slice {
				fieldValue.Set(reflect.ValueOf(value))
			} else if fieldValue.Type().Elem().Kind() == reflect.Struct && val.Len() > 0 && val.Index(0).Elem().Kind() == reflect.Map {
				slice := reflect.MakeSlice(fieldValue.Type(), val.Len(), val.Len())
				for j := 0; j < val.Len(); j++ {
					elem := reflect.New(fieldValue.Type().Elem()).Interface()
					for k, v := range val.Index(j).Interface().(map[string]interface{}) {
						if err := setField(elem, k, v); err != nil {
							return err
						}
					}
					slice.Index(j).Set(reflect.ValueOf(elem).Elem())
				}
				fieldValue.Set(slice)
			}
		} else {
			return errors.New("type mismatch for field " + name)
		}
		return nil
	}
	return nil
}
