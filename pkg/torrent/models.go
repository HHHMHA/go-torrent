package torrent

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os"
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

func NewTorrentFromBencode(data []byte) (*TorrentFile, error) {
	var torrentFile TorrentFile
	err := bencoder.NewSimpleBencoder().Unmarshal(data, &torrentFile)
	if err != nil {
		return nil, err
	}
	return &torrentFile, nil
}

func NewTorrentFromFile(filePath string) (*TorrentFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewTorrentFromBencode(data)
}

func NewTorrentFromReader(r io.Reader) (*TorrentFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewTorrentFromBencode(data)
}
