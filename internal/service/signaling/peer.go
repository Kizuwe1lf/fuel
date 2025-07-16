package signaling

import (
	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
	"github.com/pion/webrtc/v3"
	"log"
)

type Peer struct {
	pc         *webrtc.PeerConnection
	stream     signalingv1.Signaling_SignalServer
	audioTrack *webrtc.TrackLocalStaticRTP
	roomID     string
	userID     string
	userName   string
}

func (s *Service) createPeer(stream signalingv1.Signaling_SignalServer) (*Peer, error) {
	var peer *Peer
	pc, err := webrtc.NewPeerConnection(s.config)
	if err != nil {
		log.Println("Failed to create peer connection", err)
		return peer, err
	}

	audioTrack, err := s.createAudioTrack(pc)
	if err != nil {
		log.Println("Failed to create local track:", err)
		return peer, err
	}

	peer = &Peer{
		pc:         pc,
		stream:     stream,
		audioTrack: audioTrack,
	}

	return peer, nil

}

func (s *Service) destroyPeerConnection(peer *Peer) error {
	peer.pc.Close()
	if peer != nil && peer.roomID != "" {
		s.mu.Lock()
		peers := s.rooms[peer.roomID]
		for i, p := range peers {
			if p == peer {
				s.rooms[peer.roomID] = append(peers[:i], peers[i+1:]...)
				break
			}
		}
		if len(s.rooms[peer.roomID]) == 0 {
			delete(s.rooms, peer.roomID)
		}
		s.mu.Unlock()
		log.Printf("Peer %s removed from room %s", peer.userID, peer.roomID)
	}

	return nil
}

func (s *Service) createAudioTrack(pc *webrtc.PeerConnection) (*webrtc.TrackLocalStaticRTP, error) {
	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: "audio/opus",
	}, "audio", "sfu")
	if err != nil {
		return audioTrack, err
	}

	_, err = pc.AddTrack(audioTrack)
	return audioTrack, err
}
