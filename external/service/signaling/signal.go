package signaling

import signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"

func (s *Service) Signal(stream signalingv1.Signaling_SignalServer) error {
	return s.internalsignaling.Signal(stream)
}
