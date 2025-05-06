package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"torrent/pkg/peer"
	"torrent/pkg/torrent"
)

func TestNewTorrentTask_Valid(t *testing.T) {
	// Prepare a valid torrent file
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*3), // 3 pieces
			Name:        "test.torrent",
			Length:      768,
		},
	}

	tt, err := NewTorrentTask(torrentFile)

	// Valid test case
	assert.NoError(t, err)
	assert.NotNil(t, tt)
	assert.Equal(t, StatusIdle, tt.Status)
	assert.Len(t, tt.PieceStatus, 3)
	assert.Len(t, tt.Availability, 3)
}

func TestNewTorrentTask_InvalidPiecesLength(t *testing.T) {
	// Prepare a torrent with invalid piece length
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 15), // Invalid, 15 bytes (not a multiple of 20)
			Name:        "test.torrent",
			Length:      768,
		},
	}

	tt, err := NewTorrentTask(torrentFile)

	// Check for invalid piece length
	assert.Error(t, err)
	assert.Nil(t, tt)
}

func TestAddPeer(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*3),
			Name:        "test.torrent",
			Length:      768,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	peer1 := peer.Peer{IP: "192.168.1.1", Port: 6881}
	tt.AddPeer(peer1)

	// Validate peer addition
	assert.Len(t, tt.Peers, 1)
	assert.Equal(t, peer1, tt.Peers[0])
}

func TestUpdatePieceStatus(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*3),
			Name:        "test.torrent",
			Length:      768,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	// Update a piece status
	tt.UpdatePieceStatus(1)

	// Validate piece status
	assert.True(t, tt.PieceStatus[1])
	assert.Equal(t, 1.0/3.0, tt.Progress)

	// Ensure progress is updated correctly
	tt.UpdatePieceStatus(2)
	assert.True(t, tt.PieceStatus[2])
	assert.Equal(t, 2.0/3.0, tt.Progress)

	// Ensure status updates when all pieces are downloaded
	assert.Equal(t, StatusIdle, tt.Status)

	tt.UpdatePieceStatus(0)
	assert.True(t, tt.PieceStatus[0])
	assert.Equal(t, 1.0, tt.Progress)
	assert.Equal(t, StatusCompleted, tt.Status)
}

func TestGetProgress(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*4), // 4 pieces
			Name:        "test.torrent",
			Length:      1024,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	// Update some pieces
	tt.UpdatePieceStatus(1)
	tt.UpdatePieceStatus(3)

	// Check the progress
	progress := tt.GetProgress()
	assert.Equal(t, 2.0/4.0, progress)
}

func TestSetStatus(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*2), // 2 pieces
			Name:        "test.torrent",
			Length:      512,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	// Set the status to downloading
	tt.SetStatus(StatusDownloading)
	assert.Equal(t, StatusDownloading, tt.Status)

	// Set the status to paused
	tt.SetStatus(StatusPaused)
	assert.Equal(t, StatusPaused, tt.Status)
}

func TestMultiplePeerUpdates(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*4), // 4 pieces
			Name:        "test.torrent",
			Length:      1024,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	// Add multiple peers
	peer1 := peer.Peer{IP: "192.168.1.1", Port: 6881}
	peer2 := peer.Peer{IP: "192.168.1.2", Port: 6882}

	tt.AddPeer(peer1)
	tt.AddPeer(peer2)

	// Validate peers added
	assert.Len(t, tt.Peers, 2)
	assert.Equal(t, peer1, tt.Peers[0])
	assert.Equal(t, peer2, tt.Peers[1])
}

func TestInvalidPieceStatusUpdate(t *testing.T) {
	// Prepare the torrent task
	torrentFile := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		Info: torrent.InfoDict{
			PieceLength: 256,
			Pieces:      make([]byte, 20*3), // 3 pieces
			Name:        "test.torrent",
			Length:      768,
		},
	}

	tt, err := NewTorrentTask(torrentFile)
	assert.NoError(t, err)

	// Try to update an invalid piece index (e.g., negative index)
	tt.UpdatePieceStatus(-1)
	tt.UpdatePieceStatus(5)

	// Ensure piece status hasn't changed
	assert.Equal(t, 0.0, tt.Progress)
	assert.Equal(t, StatusIdle, tt.Status)
}
