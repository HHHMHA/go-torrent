package engine

import (
	"fmt"
	"sync"
	"torrent/pkg/peer"
	"torrent/pkg/torrent"
)

const PieceHashLength = 20 // SHA-1 hash length in bytes

type TorrentStatus int

const (
	StatusIdle TorrentStatus = iota
	StatusDownloading
	StatusPaused
	StatusCompleted
	StatusError
)

func (s TorrentStatus) String() string {
	return [...]string{"Idle", "Downloading", "Paused", "Completed", "Error"}[s]
}

type TorrentTask struct {
	Torrent      *torrent.TorrentFile
	Peers        []peer.Peer
	PieceStatus  []bool // true if the piece is downloaded
	Availability []int  // number of peers that have each piece
	Status       TorrentStatus
	Progress     float64
	Downloaded   int64 // bytes downloaded
	Uploaded     int64 // bytes uploaded

	mu sync.RWMutex // protects access to mutable fields
}

func NewTorrentTask(torrent *torrent.TorrentFile) (*TorrentTask, error) {
	if len(torrent.Info.Pieces)%PieceHashLength != 0 {
		return nil, fmt.Errorf("invalid pieces length: not a multiple of %d", PieceHashLength)
	}

	numPieces := len(torrent.Info.Pieces) / PieceHashLength
	return &TorrentTask{
		Torrent:      torrent,
		Peers:        []peer.Peer{},
		PieceStatus:  make([]bool, numPieces),
		Availability: make([]int, numPieces),
		Status:       StatusIdle,
	}, nil
}

func (tt *TorrentTask) AddPeer(peer peer.Peer) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.Peers = append(tt.Peers, peer)
}

func (tt *TorrentTask) UpdatePieceStatus(index int) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	if index < 0 || index >= len(tt.PieceStatus) || tt.PieceStatus[index] {
		return
	}
	tt.PieceStatus[index] = true
	completed := 0
	for _, downloaded := range tt.PieceStatus {
		if downloaded {
			completed++
		}
	}
	tt.Progress = float64(completed) / float64(len(tt.PieceStatus))
	if completed == len(tt.PieceStatus) {
		tt.Status = StatusCompleted
	}
}

func (tt *TorrentTask) GetProgress() float64 {
	tt.mu.RLock()
	defer tt.mu.RUnlock()
	return tt.Progress
}

func (tt *TorrentTask) SetStatus(status TorrentStatus) {
	tt.mu.Lock()
	defer tt.mu.Unlock()
	tt.Status = status
}
