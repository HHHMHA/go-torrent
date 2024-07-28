package models

import (
	"bytes"
	"crypto/sha1"
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
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	Files       []File `bencode:"files"`
}

type File struct {
	Length int      `bencode:"length"`
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
