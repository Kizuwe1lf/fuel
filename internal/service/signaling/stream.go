package signaling

import (
	"io"
	"log"

	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
	v1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
	"github.com/pion/webrtc/v3"
)

func (s *Service) HandleStream(peer *Peer) {
	for {
		in, err := peer.stream.Recv()
		if err == io.EOF || err != nil {
			log.Printf("Stream closed for peer %s: %v", peer.userID, err)
			break
		}

		switch msg := in.Payload.(type) {
		case *signalingv1.SignalRequest_Join:
			s.handleJoin(peer, msg)

		case *signalingv1.SignalRequest_Offer:
			s.handleOffer(peer, msg)

		case *signalingv1.SignalRequest_Answer:
			s.handleAnswer(peer, msg)

		case *signalingv1.SignalRequest_Candidate:
			s.handleCandidate(peer, msg)
		}
	}

}

func (s *Service) handleJoin(peer *Peer, req *v1.SignalRequest_Join) {
	peer.roomID = req.Join.RoomId
	peer.userID = req.Join.UserId

	s.mu.Lock()
	s.rooms[peer.roomID] = append(s.rooms[peer.roomID], peer)
	s.mu.Unlock()

	log.Printf("Peer %s joined room %s", peer.userID, peer.roomID)
}

func (s *Service) handleOffer(peer *Peer, req *v1.SignalRequest_Offer) {

	offer := webrtc.SessionDescription{
		Type: webrtc.NewSDPType(req.Offer.Type),
		SDP:  req.Offer.Sdp,
	}

	if err := peer.pc.SetRemoteDescription(offer); err != nil {
		log.Println("SetRemoteDescription failed:", err)
		return
	}

	answer, err := peer.pc.CreateAnswer(nil)
	if err != nil {
		log.Println("CreateAnswer failed:", err)
		return
	}

	if err := peer.pc.SetLocalDescription(answer); err != nil {
		log.Println("SetLocalDescription failed:", err)
		return
	}

	err = peer.stream.Send(&signalingv1.SignalResponse{
		Payload: &signalingv1.SignalResponse_Answer{
			Answer: &signalingv1.SessionDescription{
				Sdp:  answer.SDP,
				Type: answer.Type.String(),
			},
		},
	})
	if err != nil {
		log.Println("Failed to send answer:", err)
	}

}

func (s *Service) handleAnswer(peer *Peer, req *v1.SignalRequest_Answer) {
	answer := webrtc.SessionDescription{
		Type: webrtc.NewSDPType(req.Answer.Type),
		SDP:  req.Answer.Sdp,
	}

	if err := peer.pc.SetRemoteDescription(answer); err != nil {
		log.Println("SetRemoteDescription (answer) failed:", err)
	}
}

func (s *Service) handleCandidate(peer *Peer, req *v1.SignalRequest_Candidate) {
	sdpMLineIndex := uint16(req.Candidate.SdpMLineIndex)
	err := peer.pc.AddICECandidate(webrtc.ICECandidateInit{
		Candidate:     req.Candidate.Candidate,
		SDPMid:        &req.Candidate.SdpMid,
		SDPMLineIndex: &sdpMLineIndex,
	})
	if err != nil {
		log.Printf("Failed to add ICE candidate: %v", err)
	}
}
