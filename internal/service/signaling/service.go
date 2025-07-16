package signaling

import (
	"github.com/pion/webrtc/v3"
	"sync"
)

type Service struct {
	config webrtc.Configuration
	rooms  map[string][]*Peer
	mu     sync.RWMutex
}

func NewService() *Service {
	return &Service{
		rooms: make(map[string][]*Peer),
		config: webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"stun:stun.l.google.com:19302"},
				},
			},
		},
	}
}
