package peer

import "time"

type Peer struct {
	IP         string
	Port       int
	Choked     bool
	Interested bool
	LastSeen   time.Time
}
