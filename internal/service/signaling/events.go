package signaling

import (
	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"io"
	"log"
)

func (s *Service) onICECandidate(peer *Peer) {
	peer.pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		candidate := c.ToJSON()

		signalCandidate := &signalingv1.ICECandidate{
			Candidate:     candidate.Candidate,
			SdpMid:        *candidate.SDPMid,
			SdpMLineIndex: int32(*candidate.SDPMLineIndex),
		}

		err := peer.stream.Send(&signalingv1.SignalResponse{
			Payload: &signalingv1.SignalResponse_Candidate{
				Candidate: signalCandidate,
			},
		})
		if err != nil {
			log.Println("Send candidate error:", err)
		}

	})
}

func (s *Service) onTrack(peer *Peer) {
	peer.pc.OnTrack(func(track *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		log.Printf("Received track from %s: %s %s", peer.userID, track.Kind(), track.ID())

		go func() {
			rtpBuf := make([]byte, 1500)
			for {
				n, _, readErr := track.Read(rtpBuf)
				if readErr != nil {
					if readErr != io.EOF {
						log.Println("Read error:", readErr)
					}
					return
				}

				packet := &rtp.Packet{}
				if err := packet.Unmarshal(rtpBuf[:n]); err != nil {
					continue
				}

				// todo implement pub/sub syste mfor horizontal scaling
				s.mu.RLock()
				roomPeers := s.rooms[peer.roomID]
				for _, p := range roomPeers {
					if p != peer && p.audioTrack != nil {
						if err := p.audioTrack.WriteRTP(packet); err != nil {
							log.Printf("Failed to write RTP to peer %s: %v", p.userID, err)
						}
					}
				}
				s.mu.RUnlock()
			}
		}()
	})
}

func (s *Service) onConnectionStateChange(peer *Peer) {
	peer.pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Peer %s connection state: %s", peer.userID, state.String())
	})

}
