package signaling

import (
	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
)

func (s *Service) Signal(stream signalingv1.Signaling_SignalServer) error {
	peer, err := s.createPeer(stream)
	if err != nil {
		return err
	}

	//todo destroy connection

	s.onICECandidate(peer)
	s.onTrack(peer)
	s.onConnectionStateChange(peer)

	s.HandleStream(peer)

	return nil
}
