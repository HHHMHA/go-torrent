package peer

import (
	"time"
	"torrent/config"
)

type Peer struct {
	IP         string
	Port       int
	Choked     bool
	Interested bool
	LastSeen   time.Time
	ID         string
}

func GetPeerID(config *config.Config) string {
	peerID := config.PeerID
	return peerID
}
